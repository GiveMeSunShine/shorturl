package main

import (
	"fmt"
	"shorturl/shortid"
)

func main() {
	for i:=0;i<=10;i++{
		code := shortid.ShortIdFromCode()
		snowflake := shortid.ShordIdFromSnowflake()
		fmt.Println(code,":",snowflake)
	}

}
