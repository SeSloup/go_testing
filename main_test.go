package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, r := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", r.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, r.status, response.Code)
		assert.Equal(t, r.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

// tasks
const cityTest = "tula"

// TestCafeCount() — проверяет работу сервера при разных значениях параметра count
func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int // передаваемое значение count
		want  int // ожидаемое количество кафе в ответе
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList[cityTest])},
		//{0,1},
	}

	for _, r := range requests {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?count="+strconv.Itoa(r.count)+"&city="+cityTest, nil) // отправляем запрос с передаваемым количеством

		handler.ServeHTTP(res, req)
		require.Equal(t, http.StatusOK, res.Code) //запрос успешно обработан

		body := strings.TrimSpace(res.Body.String()) // ответ на запрошенное количество кафе
		var cafes []string
		if body != "" {
			cafes = strings.Split(body, ",")
		}

		assert.Equal(t, r.want, len(cafes)) //проверяем результат и ожидаание
	}
}

// TestCafeSearch() — проверяет результат поиска кафе по указанной подстроке в параметре search
func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string // передаваемое значение search
		wantCount int    // ожидаемое количество кафе по предметам из search
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
		//		{"test",2},
		//		{"фасоль", 1},
	}

	for _, r := range requests {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&search="+r.search, nil) // запрос на поиск предметов

		handler.ServeHTTP(res, req)
		require.Equal(t, http.StatusOK, res.Code) //запрос успешно обработан

		body := strings.TrimSpace(res.Body.String()) // получаем тело ответа
		var cafes []string
		if body != "" {
			cafes = strings.Split(body, ",")
		}

		assert.Equal(t, r.wantCount, len(cafes)) //проверяем результат и ожидаание
		// проверяем, что полученные в ответе кафе точно содержат переданную в search строку. (ставим эту проверку после проверки на результат/ожидание, чтобы не нарушить логику)
		for _, c := range cafes {
			assert.True(t, strings.Contains(strings.ToLower(c), strings.ToLower(r.search)))
		}

	}
}
