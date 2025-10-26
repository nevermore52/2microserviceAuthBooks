package str
import (
	"time"  
)
type Book struct {
	Id 			int    `db:"id"`
	Title 		string `db:"title"`
	Author		string `db:"author"`
	Pages 		int	   `db:"pages"`
	Readed 		bool   `db:"readed"`
	
	Timeadd 	time.Time	`db:"timeadd"`
	Timereaded	*time.Time	`db:"timeread"`
}