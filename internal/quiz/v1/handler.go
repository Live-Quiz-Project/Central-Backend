package v1

import (
	"log"
	"net/http"

	"github.com/Live-Quiz-Project/Backend/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

// ---------- Quiz related handlers ---------- //
func (h *Handler) CreateQuiz(c *gin.Context) {

	var req CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	res, er := h.Service.CreateQuiz(c.Request.Context(), &req, userID)
	if er != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": er.Error()})
		return
	}

	var qpResID *uuid.UUID = nil
	var qphResID *uuid.UUID = nil

	for _, q := range req.Questions {
		var qRes *CreateQuestionResponse
		var qpRes *CreateQuestionPoolResponse

		if q.Type == util.Pool {
			qpRes, err = h.Service.CreateQuestionPool(c.Request.Context(), &q, res.ID, res.QuizHistoryID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			log.Println("Create Question Pool PASS")
			qpResID = &qpRes.ID
			qphResID = &qpRes.QuestionPoolHistoryID
		}

		if q.IsInPool == true {
			qRes, err = h.Service.CreateQuestion(c.Request.Context(), &q, res.ID, res.QuizHistoryID, qpResID, qphResID, userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			qRes, err = h.Service.CreateQuestion(c.Request.Context(), &q, res.ID, res.QuizHistoryID, nil, nil, userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		for _, qt := range qRes.Options {
			if qst, ok := qt.(map[string]any); ok {
				if qRes.Type == util.Choice || qRes.Type == util.TrueFalse {
					_, err := h.Service.CreateChoiceOption(c.Request.Context(), &CreateChoiceOptionRequest{
						Order:   int(qst["order"].(float64)),
						Content: qst["content"].(string),
						Mark:    int(qst["mark"].(float64)),
						Color:   qst["color"].(string),
						Correct: qst["correct"].(bool),
					}, qRes.ID, qRes.QuestionHistoryID, userID)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
				} else if qRes.Type == util.ShortText || qRes.Type == util.Paragraph {
					_, err := h.Service.CreateTextOption(c.Request.Context(), &CreateTextOptionRequest{
						Order:         int(qst["order"].(float64)),
						Content:       qst["content"].(string),
						Mark:          int(qst["mark"].(float64)),
						CaseSensitive: qst["case_sensitive"].(bool),
					}, qRes.ID, qRes.QuestionHistoryID, userID)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
				}
			}
		}

		res.Questions = append(res.Questions, QuestionResponse{
			Question: Question{
				ID:             qRes.ID,
				QuizID:         qRes.QuizID,
				QuestionPoolID: qRes.QuestionPoolID,
				Type:           qRes.Type,
				Order:          qRes.Order,
				Content:        qRes.Content,
				Note:           qRes.Note,
				Media:          qRes.Media,
				UseTemplate:    qRes.UseTemplate,
				TimeLimit:      qRes.TimeLimit,
				HaveTimeFactor: qRes.HaveTimeFactor,
				TimeFactor:     qRes.TimeFactor,
				FontSize:       qRes.FontSize,
				LayoutIdx:      qRes.LayoutIdx,
				SelectUpTo:     qRes.SelectUpTo,
				CreatedAt:      qRes.CreatedAt,
				UpdatedAt:      qRes.UpdatedAt,
				DeletedAt:      qRes.DeletedAt,
			},
			Options: qRes.Options,
		})
	}

	c.JSON(http.StatusCreated, &QuizResponse{
		Quiz: Quiz{
			ID:             res.ID,
			CreatorID:      res.CreatorID,
			Title:          res.Title,
			Description:    res.Description,
			CoverImage:     res.CoverImage,
			Visibility:     res.Visibility,
			TimeLimit:      res.TimeLimit,
			HaveTimeFactor: res.HaveTimeFactor,
			TimeFactor:     res.TimeFactor,
			FontSize:       res.FontSize,
			Mark:           res.Mark,
			SelectUpTo:     res.SelectUpTo,
			CaseSensitive:  res.CaseSensitive,
			CreatedAt:      res.CreatedAt,
			UpdatedAt:      res.UpdatedAt,
			DeletedAt:      res.DeletedAt,
		},
		Questions: res.Questions,
	})
}

func (h *Handler) GetQuizzes(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	res, err := h.Service.GetQuizzes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println(res)

	r := make([]QuizResponse, 0)

	for _, q := range res {
		qpRes, err := h.Service.GetQuestionPoolsByQuizID(c.Request.Context(), q.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		// Put All Question Pools Data in Questions Response
		for _, qpr := range qpRes {
			q.Questions = append(q.Questions, QuestionResponse{
				Question: Question{
					ID:             qpr.ID,
					QuizID:         qpr.QuizID,
					Type:						"POOL",
					Order:          qpr.Order,
					Content:        qpr.Content,
					Note:           qpr.Note,
					Media:          qpr.Media,
					TimeLimit:      qpr.TimeLimit,
					HaveTimeFactor: qpr.HaveTimeFactor,
					TimeFactor:     qpr.TimeFactor,
					FontSize:       qpr.FontSize,
					CreatedAt:      qpr.CreatedAt,
					UpdatedAt:      qpr.UpdatedAt,
					DeletedAt:      qpr.DeletedAt,
				},
			})
		}

		qRes, err := h.Service.GetQuestionsByQuizID(c.Request.Context(), q.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		for _, qr := range qRes {

			if qr.Type == util.Choice || qr.Type == util.TrueFalse {
				ocRes, err := h.Service.GetChoiceOptionsByQuestionID(c.Request.Context(), qr.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}

				var oc []any
				for _, ocr := range ocRes {
					oc = append(oc, ChoiceOptioneResponse{
						ChoiceOption: ChoiceOption{
							ID:         ocr.ID,
							QuestionID: ocr.QuestionID,
							Order:      ocr.Order,
							Content:    ocr.Content,
							Mark:       ocr.Mark,
							Color:      ocr.Color,
							Correct:    ocr.Correct,
							CreatedAt:  ocr.CreatedAt,
							UpdatedAt:  ocr.UpdatedAt,
							DeletedAt:  ocr.DeletedAt,
						},
					})
				}

				q.Questions = append(q.Questions, QuestionResponse{
					Question: Question{
						ID:             qr.ID,
						QuizID:         qr.QuizID,
						QuestionPoolID: qr.QuestionPoolID,
						Type:           qr.Type,
						Order:          qr.Order,
						Content:        qr.Content,
						Note:           qr.Note,
						Media:          qr.Media,
						UseTemplate:    qr.UseTemplate,
						TimeLimit:      qr.TimeLimit,
						HaveTimeFactor: qr.HaveTimeFactor,
						TimeFactor:     qr.TimeFactor,
						FontSize:       qr.FontSize,
						LayoutIdx:      qr.LayoutIdx,
						SelectUpTo:     qr.SelectUpTo,
						CreatedAt:      qr.CreatedAt,
						UpdatedAt:      qr.UpdatedAt,
					},

					Options: oc,
				})
			} else if qr.Type == util.ShortText || qr.Type == util.Paragraph {
				otRes, err := h.Service.GetTextOptionsByQuestionID(c.Request.Context(), qr.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}

				var ot []any
				for _, otr := range otRes {
					ot = append(ot, TextOptionResponse{
						TextOption: TextOption{
							ID:            otr.ID,
							QuestionID:    otr.QuestionID,
							Order:         otr.Order,
							Mark:          otr.Mark,
							CaseSensitive: otr.CaseSensitive,
							Content:       otr.Content,
							CreatedAt:     otr.CreatedAt,
							UpdatedAt:     otr.UpdatedAt,
							DeletedAt:     otr.DeletedAt,
						},
					})
				}

				q.Questions = append(q.Questions, QuestionResponse{
					Question: Question{
						ID:             qr.ID,
						QuizID:         qr.QuizID,
						QuestionPoolID: qr.QuestionPoolID,
						Type:           qr.Type,
						Order:          qr.Order,
						Content:        qr.Content,
						Note:           qr.Note,
						Media:          qr.Media,
						UseTemplate:    qr.UseTemplate,
						TimeLimit:      qr.TimeLimit,
						HaveTimeFactor: qr.HaveTimeFactor,
						TimeFactor:     qr.TimeFactor,
						FontSize:       qr.FontSize,
						LayoutIdx:      qr.LayoutIdx,
						SelectUpTo:     qr.SelectUpTo,
						CreatedAt:      qr.CreatedAt,
						UpdatedAt:      qr.UpdatedAt,
					},
					Options: ot,
				})
			}
		}

		r = append(r, q)
	}

	c.JSON(http.StatusOK, r)
}

func (h *Handler) GetQuizByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	res, err := h.Service.GetQuizByID(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	qpRes, err := h.Service.GetQuestionPoolsByQuizID(c.Request.Context(), res.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		// Put All Question Pools Data in Questions Response
		for _, qpr := range qpRes {
			res.Questions = append(res.Questions, QuestionResponse{
				Question: Question{
					ID:             qpr.ID,
					QuizID:         qpr.QuizID,
					Type:						"POOL",
					Order:          qpr.Order,
					Content:        qpr.Content,
					Note:           qpr.Note,
					Media:          qpr.Media,
					TimeLimit:      qpr.TimeLimit,
					HaveTimeFactor: qpr.HaveTimeFactor,
					TimeFactor:     qpr.TimeFactor,
					FontSize:       qpr.FontSize,
					CreatedAt:      qpr.CreatedAt,
					UpdatedAt:      qpr.UpdatedAt,
					DeletedAt:      qpr.DeletedAt,
				},
			})
		}

	qRes, err := h.Service.GetQuestionsByQuizID(c.Request.Context(), res.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, qr := range qRes {
		if qr.Type == util.Choice || qr.Type == util.TrueFalse {
			ocRes, err := h.Service.GetChoiceOptionsByQuestionID(c.Request.Context(), qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var oc []any
			for _, ocr := range ocRes {
				oc = append(oc, ChoiceOptioneResponse{
					ChoiceOption: ChoiceOption{
						ID:         ocr.ID,
						QuestionID: ocr.QuestionID,
						Order:      ocr.Order,
						Content:    ocr.Content,
						Mark:       ocr.Mark,
						Color:      ocr.Color,
						Correct:    ocr.Correct,
						CreatedAt:  ocr.CreatedAt,
						UpdatedAt:  ocr.UpdatedAt,
						DeletedAt:  ocr.DeletedAt,
					},
				})
			}

			res.Questions = append(res.Questions, QuestionResponse{
				Question: Question{
					ID:             qr.ID,
					QuizID:         qr.QuizID,
					QuestionPoolID: qr.QuestionPoolID,
					Type:           qr.Type,
					Order:          qr.Order,
					Content:        qr.Content,
					Note:           qr.Note,
					Media:          qr.Media,
					UseTemplate:    qr.UseTemplate,
					TimeLimit:      qr.TimeLimit,
					HaveTimeFactor: qr.HaveTimeFactor,
					TimeFactor:     qr.TimeFactor,
					FontSize:       qr.FontSize,
					LayoutIdx:      qr.LayoutIdx,
					SelectUpTo:     qr.SelectUpTo,
					CreatedAt:      qr.CreatedAt,
					UpdatedAt:      qr.UpdatedAt,
				},
				Options: oc,
			})
		} else if qr.Type == util.ShortText || qr.Type == util.Paragraph {
			otRes, err := h.Service.GetTextOptionsByQuestionID(c.Request.Context(), qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var ot []any
			for _, otr := range otRes {
				ot = append(ot, TextOptionResponse{
					TextOption: TextOption{
						ID:            otr.ID,
						QuestionID:    otr.QuestionID,
						Order:         otr.Order,
						Content:       otr.Content,
						Mark:          otr.Mark,
						CaseSensitive: otr.CaseSensitive,
						CreatedAt:     otr.CreatedAt,
						UpdatedAt:     otr.UpdatedAt,
						DeletedAt:     otr.DeletedAt,
					},
				})
			}

			res.Questions = append(res.Questions, QuestionResponse{
				Question: Question{
					ID:             qr.ID,
					QuizID:         qr.QuizID,
					QuestionPoolID: qr.QuestionPoolID,
					Type:           qr.Type,
					Order:          qr.Order,
					Content:        qr.Content,
					Note:           qr.Note,
					Media:          qr.Media,
					UseTemplate:    qr.UseTemplate,
					TimeLimit:      qr.TimeLimit,
					HaveTimeFactor: qr.HaveTimeFactor,
					TimeFactor:     qr.TimeFactor,
					FontSize:       qr.FontSize,
					LayoutIdx:      qr.LayoutIdx,
					SelectUpTo:     qr.SelectUpTo,
					CreatedAt:      qr.CreatedAt,
					UpdatedAt:      qr.UpdatedAt,
				},
				Options: ot,
			})
		}
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateQuiz(c *gin.Context) {
	var req UpdateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}
	quizID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	res, err := h.Service.UpdateQuiz(c.Request.Context(), &req, quizID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) DeleteQuiz(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}
	quizID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	er := h.Service.DeleteQuiz(c.Request.Context(), quizID, userID)
	if er != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": er.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted",
	})
}

func (h *Handler) RestoreQuiz(c *gin.Context) {}
