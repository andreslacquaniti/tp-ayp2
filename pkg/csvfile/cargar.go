package csvfile

import (
	"encoding/csv"
	"errors"
	"os"
	"reflect"
	"strconv"
	"time"
)

// CargaCSV lee un CSV en filename usando el delimiter dado y devuelve []T.
func CargaCSV[T any](filename string) ([]T, error) {
	file, err := os.Open("data/" + filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = ';'

	// Leer encabezados
	headers, err := r.Read()
	if err != nil {
		return nil, err
	}
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var result []T
	for _, row := range rows {
		var item T
		v := reflect.ValueOf(&item).Elem()
		t := v.Type()

		for i, colName := range headers {
			for j := 0; j < t.NumField(); j++ {
				field := t.Field(j)
				tag := field.Tag.Get("csv")
				if tag == colName {
					fv := v.Field(j)
					switch fv.Kind() {
					case reflect.String:
						fv.SetString(row[i])
					case reflect.Int:
						if n, err := strconv.Atoi(row[i]); err == nil {
							fv.SetInt(int64(n))
						}
					case reflect.Float64:
						if f, err := strconv.ParseFloat(row[i], 64); err == nil {
							fv.SetFloat(f)
						}
					case reflect.Bool:
						if b, err := strconv.ParseBool(row[i]); err == nil {
							fv.SetBool(b)
						}
					default:
						// time.Time
						if field.Type == reflect.TypeOf(time.Time{}) {
							if tval, err := time.Parse("2006-01-02 15:04:05", row[i]); err == nil {
								fv.Set(reflect.ValueOf(tval))
							} else {
								return nil, err
							}
						} else {
							return nil, errors.New("tipo no soportado: " + field.Name)
						}
					}
				}
			}
		}
		result = append(result, item)
	}

	return result, nil
}
