package main

import (
	"fmt"
	"time"

	"github.com/dwiyanrp/go-scheduler"
	"github.com/gin-gonic/gin"
)

func main() {
	s := scheduler.NewScheduler()
	r := gin.New()

	r.GET("/start", func(c *gin.Context) {
		runAt := time.Now().Add(10 * time.Second)
		msg := c.Query("msg")
		taskID, _ := s.RunAt(runAt, PrintMessage, msg)
		c.String(200, fmt.Sprintf("Task %v scheduled", taskID))
	})

	r.GET("/every", func(c *gin.Context) {
		runEvery := 5 * time.Second
		runUntil := time.Now().Add(60 * time.Second)
		msg := c.Query("msg")
		taskID, _ := s.RunEvery(runEvery, runUntil, PrintMessage, msg)
		c.String(200, fmt.Sprintf("Task %v scheduled", taskID))
	})

	r.GET("/stop/:id", func(c *gin.Context) {
		taskID := c.Param("id")
		if err := s.Cancel(taskID); err != nil {
			c.String(200, fmt.Sprint(err))
		} else {
			c.String(200, fmt.Sprintf("Task %v stopped", taskID))
		}
	})

	r.GET("/reschedule/:id", func(c *gin.Context) {
		taskID := c.Param("id")
		rescheduleTime := time.Now().Add(5 * time.Second)
		if err := s.Reschedule(taskID, rescheduleTime); err != nil {
			c.String(200, fmt.Sprint(err))
		} else {
			c.String(200, fmt.Sprintf("Task %v rescheduled", taskID))
		}
	})

	r.Run()
}

func PrintMessage(msg string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), msg)
}
