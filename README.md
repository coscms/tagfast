tagfast
=======

golang：优化读取struct内的tag值（只解析一次，以后都从缓存中读取。官方的版本每次使用typ.Field(i).Tag.Get("tag1")都要解析一次，效率不高）

用法
=======
```package main
import "github.com/coscms/tagfast"

type Coscms struct {
  Url string `xorm:"not null default '' VARCHAR(255)" valid:"Requied" form_widget:"text"`
  Email string `xorm:"not null default '' VARCHAR(255)" valid:"Requied" form_widget:"text"`
}

func main(){
  m:=Coscms{}
  t := reflect.TypeOf(m)
  for i := 0; i < t.NumField(); i++ {
    widget:=tagfast.Tag(m,i,"form_widget")
  }
}```
