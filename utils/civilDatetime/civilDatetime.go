package civilDatetime

import (
	"strings"
	"time"
)

type CivilDate time.Time
type CivilTime time.Time

// 以下程式碼主要用於 golnag 日期型態問題，
// 對於只需要單純的 Y-m-d(date) or H:i:s(time) 這類型的 db type，
// 因此會需要有辦法在 response 時轉換一次格式再回傳給使用者
func (c *CivilDate) UnmarshalJSONDate(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", value) //parse date
	if err != nil {
		return err
	}
	*c = CivilDate(t) //set result using the pointer
	return nil
}

func (c CivilDate) MarshalJSONDate() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format("2006-01-02") + `"`), nil
}

func (c *CivilTime) UnmarshalJSONTime(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("15:04:05", value) //parse time
	if err != nil {
		return err
	}
	*c = CivilTime(t) //set result using the pointer
	return nil
}

func (c CivilTime) MarshalJSONTime() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format("15:04:05") + `"`), nil
}

// Another version
// func (c *CivilTime) UnmarshalJSON(data []byte) error {
// 	if string(data) == "null" {
// 		return nil
// 	}
// 	var t time.Time
// 	err := json.Unmarshal(data, &t)
// 	if err != nil {
// 		return err
// 	}
// 	*c = CivilTime(t)
// 	return nil
// }

// func (c CivilTime) MarshalJSON() ([]byte, error) {
// 	if time.Time(c).IsZero() {
// 		return []byte("123"), nil
// 	}
// 	return []byte(`"` + time.Time(c).Format("2006-01-02") + `"`), nil
// }
