package timer

import "encoding/json"

type ws struct {
	Wheels map[int]interface{} `json:"wheels"`
}

type w struct {
	Slots map[int]interface{} `json:"slots"`
}

type s struct {
	Ticks map[int]interface{} `json:"ticks"`
}

//for debugging
func (t *timingwheel) Topology() ([]byte, error) {
	ws := &ws{
		Wheels: map[int]interface{}{},
	}
	for i, wheel := range t.wheels {
		w := &w{
			Slots: map[int]interface{}{},
		}
		for j, slot := range wheel.slots {
			s := &s{
				Ticks: map[int]interface{}{},
			}
			k := 0
			slot.foreach(func(data interface{}) error {
				t := data.(*tick)
				s.Ticks[k] = t.duration
				k++
				return nil
			})
			w.Slots[j] = s
		}
		ws.Wheels[i] = w
	}
	return json.Marshal(ws)
}
