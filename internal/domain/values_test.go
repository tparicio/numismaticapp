package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/antonioparicio/numismaticapp/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestYear(t *testing.T) {
	t.Run("Valid Year", func(t *testing.T) {
		y, err := domain.NewYear(2000)
		assert.NoError(t, err)
		assert.Equal(t, 2000, y.Int())
	})

	t.Run("Zero Year", func(t *testing.T) {
		y, err := domain.NewYear(0)
		assert.NoError(t, err)
		assert.Equal(t, 0, y.Int())
	})

	t.Run("Invalid Year Low", func(t *testing.T) {
		_, err := domain.NewYear(-5001)
		assert.Error(t, err)
	})

	t.Run("Invalid Year High", func(t *testing.T) {
		_, err := domain.NewYear(3001)
		assert.Error(t, err)
	})

	t.Run("JSON Marshaling", func(t *testing.T) {
		y, _ := domain.NewYear(1999)
		data, err := json.Marshal(y)
		assert.NoError(t, err)
		assert.Equal(t, "1999", string(data))

		var y2 domain.Year
		err = json.Unmarshal(data, &y2)
		assert.NoError(t, err)
		assert.Equal(t, 1999, y2.Int())
	})

	t.Run("JSON Unmarshal Error", func(t *testing.T) {
		var y domain.Year
		err := json.Unmarshal([]byte(`"invalid"`), &y) // Expecting int
		assert.Error(t, err)

		err = json.Unmarshal([]byte(`999999`), &y) // Out of range
		assert.Error(t, err)
	})
}

func TestMintage(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		m, err := domain.NewMintage(100)
		assert.NoError(t, err)
		assert.Equal(t, int64(100), m.Int64())
	})

	t.Run("Negative", func(t *testing.T) {
		_, err := domain.NewMintage(-1)
		assert.Error(t, err)
	})

	t.Run("JSON", func(t *testing.T) {
		m, _ := domain.NewMintage(500)
		data, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, "500", string(data))

		var m2 domain.Mintage
		err = json.Unmarshal(data, &m2)
		assert.NoError(t, err)
		assert.Equal(t, int64(500), m2.Int64())
	})

	t.Run("JSON Unmarshal Error", func(t *testing.T) {
		var m domain.Mintage
		err := json.Unmarshal([]byte(`"abc"`), &m) // Expecting int64
		assert.Error(t, err)

		err = json.Unmarshal([]byte(`-5`), &m) // Negative
		assert.Error(t, err)
	})
}

func TestKMCode(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		k, err := domain.NewKMCode("KM# 123")
		assert.NoError(t, err)
		assert.Equal(t, "KM# 123", k.String())
	})

	t.Run("Empty", func(t *testing.T) {
		k, err := domain.NewKMCode("")
		assert.NoError(t, err)
		assert.Equal(t, "", k.String())
	})

	t.Run("JSON", func(t *testing.T) {
		k, _ := domain.NewKMCode("Y# 5")
		data, err := json.Marshal(k)
		assert.NoError(t, err)
		assert.Equal(t, `"Y# 5"`, string(data))

		var k2 domain.KMCode
		err = json.Unmarshal(data, &k2)
		assert.NoError(t, err)
		assert.Equal(t, "Y# 5", k2.String())
	})

	t.Run("JSON Unmarshal Error", func(t *testing.T) {
		var k domain.KMCode
		err := json.Unmarshal([]byte(`123`), &k) // Expecting string
		assert.Error(t, err)
	})
}

func TestGrade(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		g, err := domain.NewGrade("EBC")
		assert.NoError(t, err)
		assert.Equal(t, "EBC", g.String())
	})

	t.Run("JSON", func(t *testing.T) {
		g, _ := domain.NewGrade("UNC")
		data, err := json.Marshal(g)
		assert.NoError(t, err)
		assert.Equal(t, `"UNC"`, string(data))

		var g2 domain.Grade
		err = json.Unmarshal(data, &g2)
		assert.NoError(t, err)
		assert.Equal(t, "UNC", g2.String())
	})

	t.Run("JSON Unmarshal Error", func(t *testing.T) {
		var g domain.Grade
		err := json.Unmarshal([]byte(`123`), &g) // Expecting string
		assert.Error(t, err)
	})
}
