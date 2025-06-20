package csvfile

import (
	"encoding/csv"
	"errors"
	"os"
	"reflect"
	"strconv"
	"time"
)

// GuardarCSV escribe en filename (trunca/crea) el slice registros con el delimiter dado.
func GuardarCSV[T any](filename string, registros []T) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	w.Comma = ';'

	// Cabecera: extraemos las etiquetas csv en orden de campos
	var headers []string
	var fields []reflect.StructField

	typ := reflect.TypeOf((*T)(nil)).Elem()
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if tag := f.Tag.Get("csv"); tag != "" {
			headers = append(headers, tag)
			fields = append(fields, f)
		}
	}
	if err := w.Write(headers); err != nil {
		return err
	}

	// Filas: por cada registro, recorremos los campos
	for _, rec := range registros {
		v := reflect.ValueOf(rec)
		var row []string

		for _, f := range fields {
			fv := v.FieldByName(f.Name)
			switch fv.Kind() {
			case reflect.String:
				row = append(row, fv.String())
			case reflect.Int:
				row = append(row, strconv.FormatInt(fv.Int(), 10))
			case reflect.Float64:
				row = append(row, strconv.FormatFloat(fv.Float(), 'f', -1, 64))
			case reflect.Bool:
				row = append(row, strconv.FormatBool(fv.Bool()))
			default:
				if f.Type == reflect.TypeOf(time.Time{}) {
					tval := fv.Interface().(time.Time)
					row = append(row, tval.Format("2006-01-02 15:04:05"))
				} else {
					return errors.New("tipo no soportado al escribir: " + f.Name)
				}
			}
		}

		if err := w.Write(row); err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}
