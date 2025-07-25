package generator

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// Параметры генерации пароля / Parameters for password generation

const (
	LengthMin = 4
	LengthMax = 35
	ups       = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lows      = "abcdefghijklmnopqrstuvwxyz"
	digs      = "0123456789"
	specSymbs = "!@#№$;%^:&?*()-_=+[]{}<>.,/|`~"
	similars  = "il1O0"
)

var (
	ErrLengthOutOfRange   = errors.New("length_out_of_range")
	ErrNoCategorySelected = errors.New("no_category_selected")
)

type FlagsSet struct {
	Upper, Lower, Digits, SpecSymbols, ExcludeSimilar bool
}

func (f FlagsSet) HasAny() bool {
	return f.Upper || f.Lower || f.Digits || f.SpecSymbols
}

// Generate создаёт пароль длины LengthMin–LengthMax, минимум по одному символу из каждой выбранной категории
// Generate creates a password of lenghtMin-LenghtMax, at least one symbol from each selected category

func Generate(length int, flags FlagsSet) (string, error) {
	if length < LengthMin || length > LengthMax {
		return "", ErrLengthOutOfRange
	}

	var pool []rune
	var required [][]rune

	if flags.Upper {
		runes := []rune(ups)
		required = append(required, runes)
		pool = append(pool, runes...)
	}
	if flags.Lower {
		runes := []rune(lows)
		required = append(required, runes)
		pool = append(pool, runes...)
	}
	if flags.Digits {
		runes := []rune(digs)
		required = append(required, runes)
		pool = append(pool, runes...)
	}
	if flags.SpecSymbols {
		runes := []rune(specSymbs)
		required = append(required, runes)
		pool = append(pool, runes...)
	}

	if len(required) == 0 {
		return "", ErrNoCategorySelected
	}
	if len(required) > length {
		return "", ErrLengthOutOfRange
	}

	if flags.ExcludeSimilar {

		// Фильтруются похожие символы из pool и required / Filter similar characters from pool and required

		filteredPool := make([]rune, 0, len(pool))
		for _, r := range pool {
			if !containsRune(similars, r) {
				filteredPool = append(filteredPool, r)
			}
		}
		pool = filteredPool

		for i, cat := range required {
			filteredCat := make([]rune, 0, len(cat))
			for _, r := range cat {
				if !containsRune(similars, r) {
					filteredCat = append(filteredCat, r)
				}
			}
			required[i] = filteredCat
		}
	}

	// Собирается пароль: сначала по одному символу из каждой категории / Build password: first one symbol from each category

	password := make([]rune, length)
	for i, cat := range required {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(cat))))
		if err != nil {
			return "", err
		}
		password[i] = cat[idx.Int64()]
	}

	// Заполняются оставшиеся позиции случайными символами из pool / Fill the rest of the positions with random symbols from pool

	for i := len(required); i < length; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(pool))))
		if err != nil {
			return "", err
		}
		password[i] = pool[idx.Int64()]
	}

	// Перемешивается массив рун / Shuffle the array of runes

	mix(password)
	return string(password), nil
}

func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}

func mix(runes []rune) error {
	n := len(runes)
	for i := n - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return err
		}
		k := int(j.Int64())
		runes[i], runes[k] = runes[k], runes[i]
	}
	return nil
}
