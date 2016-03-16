package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionInterviews = "interviews"
)

type Interview struct {
	ID        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title"`
	Date      time.Time     `json:"date" bson:"date" binding:"required"`
	Owner     bson.ObjectId `json:"owner" bson:"owner"`
	Vacancy   bson.ObjectId `json:"vacancy" bson:"vacancy" binding:"required"`
	Candidate bson.ObjectId `json:"candidate" bson:"candidate" binding:"required"`
	Company   bson.ObjectId `json:"company" bson:"company" binding:"required"`
}

type InterviewViewModel struct {
	ID        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string        `json:"title,omitempty" bson:"title"`
	Date      *time.Time    `json:"date,omitempty" bson:"date"`
	Owner     *User         `json:"owner,omitempty" bson:"owner"`
	Vacancy   *Vacancy      `json:"vacancy,omitempty" bson:"vacancy" binding:"required"`
	Candidate *User         `json:"candidate,omitempty" bson:"candidate" binding:"required"`
	Company   *Company      `json:"company,omitempty" bson:"company" binding:"required"`
}
