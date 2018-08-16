package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dwiyanrp/go-scheduler"
	"github.com/gin-gonic/gin"
)

func main() {
	s := scheduler.NewScheduler()
	r := gin.New()

	r.GET("/start", func(c *gin.Context) {
		taskID, _ := s.RunAt(time.Now().Add(10*time.Second), PrintMessage, c.Query("msg"))
		c.String(200, fmt.Sprintf("Task %v scheduled", taskID))
	})

	r.GET("/stop/:id", func(c *gin.Context) {
		taskID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		if err := s.Cancel(taskID); err != nil {
			c.String(200, fmt.Sprint(err))
		} else {
			c.String(200, fmt.Sprintf("Task %v stopped", taskID))
		}
	})

	r.GET("/reschedule/:id", func(c *gin.Context) {
		taskID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		if err := s.Reschedule(taskID, time.Now().Add(5*time.Second)); err != nil {
			c.String(200, fmt.Sprint(err))
		} else {
			c.String(200, fmt.Sprintf("Task %v rescheduled", taskID))
		}
	})

	r.Run()
}

func PrintMessage(msg string) {
	fmt.Println(msg)
}
