package v1

import (
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

	// Start Transaction
	tx, er := h.Service.BeginTransaction(c.Request.Context())
	if er != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": er.Error()})
		return
	}

	res, err := h.Service.CreateQuiz(c.Request.Context(), tx, &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var qpResID *uuid.UUID = nil
	var qphResID *uuid.UUID = nil

	for _, q := range req.Questions {
		var qRes *CreateQuestionResponse
		var qpRes *CreateQuestionPoolResponse

		if q.Type == util.Pool {
			qpRes, err = h.Service.CreateQuestionPool(c.Request.Context(), tx, &q, res.ID, res.QuizHistoryID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			qpResID = &qpRes.ID
			qphResID = &qpRes.QuestionPoolHistoryID

			res.Questions = append(res.Questions, QuestionResponse{
				Question: Question{
					ID:             qpRes.ID,
					QuizID:         qpRes.QuizID,
					Type:           "POOL",
					Order:          qpRes.Order,
					Content:        qpRes.Content,
					Note:           qpRes.Note,
					Media:          qpRes.Media,
					TimeLimit:      qpRes.TimeLimit,
					HaveTimeFactor: qpRes.HaveTimeFactor,
					TimeFactor:     qpRes.TimeFactor,
					FontSize:       qpRes.FontSize,
					CreatedAt:      qpRes.CreatedAt,
					UpdatedAt:      qpRes.UpdatedAt,
					DeletedAt:      qpRes.DeletedAt,
				},
			})

			continue
		} else {
			if q.IsInPool == true {
				qRes, err = h.Service.CreateQuestion(c.Request.Context(), tx, &q, res.ID, res.QuizHistoryID, qpResID, qphResID, userID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			} else {
				qRes, err = h.Service.CreateQuestion(c.Request.Context(), tx, &q, res.ID, res.QuizHistoryID, nil, nil, userID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}

			for _, qt := range qRes.Options {
				if qst, ok := qt.(map[string]any); ok {
					if qRes.Type == util.Choice || qRes.Type == util.TrueFalse {
						_, err := h.Service.CreateChoiceOption(c.Request.Context(), tx, &ChoiceOptionRequest{
							ChoiceOption: ChoiceOption{
								Order:   int(qst["order"].(float64)),
								Content: qst["content"].(string),
								Mark:    int(qst["mark"].(float64)),
								Color:   qst["color"].(string),
								Correct: qst["correct"].(bool),
							},
						}, qRes.ID, qRes.QuestionHistoryID, userID)

						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
							return
						}

					} else if qRes.Type == util.ShortText || qRes.Type == util.Paragraph {
						_, err := h.Service.CreateTextOption(c.Request.Context(), tx, &TextOptionRequest{
							TextOption: TextOption{
								Order:         int(qst["order"].(float64)),
								Content:       qst["content"].(string),
								Mark:          int(qst["mark"].(float64)),
								CaseSensitive: qst["case_sensitive"].(bool),
							},
						}, qRes.ID, qRes.QuestionHistoryID, userID)

						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
							return
						}

					} else if qRes.Type == util.Matching {
						if qst["type"].(string) != "MATCHING_ANSWER" {
							_, err := h.Service.CreateMatchingOption(c.Request.Context(), tx, &MatchingOptionRequest{
								MatchingOption: MatchingOption{
									Order:     int(qst["order"].(float64)),
									Content:   qst["content"].(string),
									Type:      qst["type"].(string),
									Eliminate: qst["eliminate"].(bool),
								},
							}, qRes.ID, qRes.QuestionHistoryID, userID)

							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

						} else {

							prompt, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c.Request.Context(), qRes.ID, int(qst["prompt"].(float64)))
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							option, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c.Request.Context(), qRes.ID, int(qst["option"].(float64)))
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							_, err = h.Service.CreateMatchingAnswer(c.Request.Context(), tx, &MatchingAnswerRequest{
								MatchingAnswer: MatchingAnswer{
									PromptID: prompt.ID,
									OptionID: option.ID,
									Mark:     int(qst["mark"].(float64)),
								},
							}, qRes.ID, qRes.QuestionHistoryID, userID)

							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}
						}
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

	er = h.Service.CommitTransaction(c.Request.Context(), tx)

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
					Type:           "POOL",
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
					oc = append(oc, ChoiceOptionResponse{
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
			} else if qr.Type == util.Matching {
				omRes, err := h.Service.GetMatchingOptionsByQuestionID(c.Request.Context(), qr.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				amRes, err := h.Service.GetMatchingAnswersByQuestionID(c.Request.Context(), qr.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				var o []any
				for _, omr := range omRes {
					o = append(o, MatchingOptionAndAnswerResponse{
						ID:         omr.ID,
						QuestionID: omr.QuestionID,
						Type:       omr.Type,
						Order:      omr.Order,
						Content:    omr.Content,
						Eliminate:  omr.Eliminate,
						PromptID:   uuid.Nil,
						OptionID:   uuid.Nil,
						Mark:       0,
						CreatedAt:  omr.CreatedAt,
						UpdatedAt:  omr.UpdatedAt,
						DeletedAt:  omr.DeletedAt,
					})
				}

				for _, amr := range amRes {
					o = append(o, MatchingOptionAndAnswerResponse{
						ID:         amr.ID,
						QuestionID: amr.QuestionID,
						Type:       "MATCHING_ANSWER",
						Order:      0,
						Content:    "",
						Eliminate:  false,
						PromptID:   amr.PromptID,
						OptionID:   amr.OptionID,
						Mark:       amr.Mark,
						CreatedAt:  amr.CreatedAt,
						UpdatedAt:  amr.UpdatedAt,
						DeletedAt:  amr.DeletedAt,
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
					Options: o,
				})
			}
		}

		r = append(r, q)
	}

	// Bubble Sort the Questions and Question Pool by Order
	for _, res := range r {
		n := len(res.Questions)
		for i := 0; i < n-1; i++ {
			for j := 0; j < n-i-1; j++ {
				if res.Questions[j].Order > res.Questions[j+1].Order {
					res.Questions[j], res.Questions[j+1] = res.Questions[j+1], res.Questions[j]
				}
			}
		}
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
				Type:           "POOL",
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
				oc = append(oc, ChoiceOptionResponse{
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
		} else if qr.Type == util.Matching {
			omRes, err := h.Service.GetMatchingOptionsByQuestionID(c.Request.Context(), qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			amRes, err := h.Service.GetMatchingAnswersByQuestionID(c.Request.Context(), qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var o []any
			for _, omr := range omRes {
				o = append(o, MatchingOptionAndAnswerResponse{
					ID:         omr.ID,
					QuestionID: omr.QuestionID,
					Type:       omr.Type,
					Order:      omr.Order,
					Content:    omr.Content,
					Eliminate:  omr.Eliminate,
					PromptID:   uuid.Nil,
					OptionID:   uuid.Nil,
					Mark:       0,
					CreatedAt:  omr.CreatedAt,
					UpdatedAt:  omr.UpdatedAt,
					DeletedAt:  omr.DeletedAt,
				})
			}

			for _, amr := range amRes {
				o = append(o, MatchingOptionAndAnswerResponse{
					ID:         amr.ID,
					QuestionID: amr.QuestionID,
					Type:       "MATCHING_ANSWER",
					Order:      0,
					Content:    "",
					Eliminate:  false,
					PromptID:   amr.PromptID,
					OptionID:   amr.OptionID,
					Mark:       amr.Mark,
					CreatedAt:  amr.CreatedAt,
					UpdatedAt:  amr.UpdatedAt,
					DeletedAt:  amr.DeletedAt,
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
				Options: o,
			})
		}
	}

	// Bubble Sort the Questions and Question Pool by Order
	n := len(res.Questions)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if res.Questions[j].Order > res.Questions[j+1].Order {
				res.Questions[j], res.Questions[j+1] = res.Questions[j+1], res.Questions[j]
			}
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

	tx, err := h.Service.BeginTransaction(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.UpdateQuiz(c.Request.Context(), tx, &req, userID, quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var qpResID *uuid.UUID = nil
	var qphResID *uuid.UUID = nil

	for _, q := range req.Questions {

		if q.ID != uuid.Nil {
			var qRes *UpdateQuestionResponse
			var qpRes *UpdateQuestionPoolResponse

			if q.Type == util.Pool {
				qpRes, err = h.Service.UpdateQuestionPool(c.Request.Context(), tx, &q, userID, q.ID, res.QuizHistoryID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}
				qpResID = &qpRes.ID
				qphResID = &qpRes.QuestionPoolHistoryID

				res.Questions = append(res.Questions, QuestionResponse{
					Question: Question{
						ID:             qpRes.ID,
						QuizID:         qpRes.QuizID,
						Type:           "POOL",
						Order:          qpRes.Order,
						Content:        qpRes.Content,
						Note:           qpRes.Note,
						Media:          qpRes.Media,
						TimeLimit:      qpRes.TimeLimit,
						HaveTimeFactor: qpRes.HaveTimeFactor,
						TimeFactor:     qpRes.TimeFactor,
						FontSize:       qpRes.FontSize,
						CreatedAt:      qpRes.CreatedAt,
						UpdatedAt:      qpRes.UpdatedAt,
						DeletedAt:      qpRes.DeletedAt,
					},
				})
				continue

			} else {
				if q.IsInPool == true {
					qRes, err = h.Service.UpdateQuestion(c.Request.Context(), tx, &q, userID, q.ID, res.QuizHistoryID, qphResID)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
				} else {
					qRes, err = h.Service.UpdateQuestion(c.Request.Context(), tx, &q, userID, q.ID, res.QuizHistoryID, nil)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
				}
			}

			// Option in Question Checking
			for _, qt := range qRes.Options {
				if qst, ok := qt.(map[string]any); ok {
					var id uuid.UUID
					var questionID uuid.UUID

					if strID, ok := qst["id"].(string); ok {
						id, _ = uuid.Parse(strID)
					}

					if strQuestionID, ok := qst["question_id"].(string); ok {
						questionID, _ = uuid.Parse(strQuestionID)
					}

					if qRes.Type == util.Choice || qRes.Type == util.TrueFalse {
						choiceReq := ChoiceOptionRequest{
							ChoiceOption: ChoiceOption{
								ID:         id,
								QuestionID: questionID,
								Order:      int(qst["order"].(float64)),
								Content:    qst["content"].(string),
								Mark:       int(qst["mark"].(float64)),
								Color:      qst["color"].(string),
								Correct:    qst["correct"].(bool),
							},
						}

						if choiceReq.ID != uuid.Nil {
							_, err := h.Service.UpdateChoiceOption(c.Request.Context(), tx, &choiceReq, userID, id, qRes.QuestionHistoryID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

						} else {
							_, err := h.Service.CreateChoiceOption(c.Request.Context(), tx, &choiceReq, qRes.ID, qRes.QuestionHistoryID, userID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}
						}

					} else if qRes.Type == util.ShortText || qRes.Type == util.Paragraph {
						textReq := TextOptionRequest{
							TextOption: TextOption{
								ID:            id,
								QuestionID:    questionID,
								Order:         int(qst["order"].(float64)),
								Content:       qst["content"].(string),
								Mark:          int(qst["mark"].(float64)),
								CaseSensitive: qst["case_sensitive"].(bool),
							},
						}

						if textReq.ID != uuid.Nil {
							_, err := h.Service.UpdateTextOption(c.Request.Context(), tx, &textReq, userID, id, qRes.QuestionHistoryID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

						} else {
							_, err := h.Service.CreateTextOption(c.Request.Context(), tx, &textReq, qRes.ID, qRes.QuestionHistoryID, userID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}
						}
					} else if qRes.Type == util.Matching {
						if qst["type"].(string) != "MATCHING_ANSWER" {
							matchingOptionReq := MatchingOptionRequest{
								MatchingOption: MatchingOption{
									ID:         id,
									QuestionID: questionID,
									Type:       qst["type"].(string),
									Order:      int(qst["order"].(float64)),
									Content:    qst["content"].(string),
									Eliminate:  qst["eliminate"].(bool),
								},
							}

							if matchingOptionReq.ID != uuid.Nil {
								_, err := h.Service.UpdateMatchingOption(c.Request.Context(), tx, &matchingOptionReq, userID, id, qRes.QuestionHistoryID)
								if err != nil {
									c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
									return
								}

							} else {
								_, err := h.Service.CreateMatchingOption(c.Request.Context(), tx, &matchingOptionReq, qRes.ID, qRes.QuestionHistoryID, userID)
								if err != nil {
									c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
									return
								}
							}
						} else {

							prompt, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c.Request.Context(), qRes.ID, int(qst["prompt"].(float64)))
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							option, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c.Request.Context(), qRes.ID, int(qst["option"].(float64)))
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							if id != uuid.Nil {
								_, err = h.Service.UpdateMatchingAnswer(c.Request.Context(), tx, &MatchingAnswerRequest{
									MatchingAnswer: MatchingAnswer{
										ID:         id,
										QuestionID: questionID,
										PromptID:   prompt.ID,
										OptionID:   option.ID,
										Mark:       int(qst["mark"].(float64)),
									},
								}, userID, id, qRes.QuestionHistoryID)
							} else {
								_, err = h.Service.CreateMatchingAnswer(c.Request.Context(), tx, &MatchingAnswerRequest{
									MatchingAnswer: MatchingAnswer{
										PromptID: prompt.ID,
										OptionID: option.ID,
										Mark:     int(qst["mark"].(float64)),
									},
								}, qRes.ID, qRes.QuestionHistoryID, userID)
							}
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

			// If the new one
		} else {
			var qRes *CreateQuestionResponse
			var qpRes *CreateQuestionPoolResponse

			if q.Type == util.Pool {
				qpRes, err = h.Service.CreateQuestionPool(c.Request.Context(), tx, &q, quizID, res.QuizHistoryID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				qpResID = &qpRes.ID
				qphResID = &qpRes.QuestionPoolHistoryID
				continue
			}

			if q.IsInPool == true {
				qRes, err = h.Service.CreateQuestion(c.Request.Context(), tx, &q, quizID, res.QuizHistoryID, qpResID, qphResID, userID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			} else {
				qRes, err = h.Service.CreateQuestion(c.Request.Context(), tx, &q, quizID, res.QuizHistoryID, nil, nil, userID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}

			for _, qt := range qRes.Options {
				if qst, ok := qt.(map[string]any); ok {
					if qRes.Type == util.Choice || qRes.Type == util.TrueFalse {
						_, err := h.Service.CreateChoiceOption(c.Request.Context(), tx, &ChoiceOptionRequest{
							ChoiceOption: ChoiceOption{
								Order:   int(qst["order"].(float64)),
								Content: qst["content"].(string),
								Mark:    int(qst["mark"].(float64)),
								Color:   qst["color"].(string),
								Correct: qst["correct"].(bool),
							},
						}, qRes.ID, qRes.QuestionHistoryID, userID)
						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
							return
						}
					} else if qRes.Type == util.ShortText || qRes.Type == util.Paragraph {
						_, err := h.Service.CreateTextOption(c.Request.Context(), tx, &TextOptionRequest{
							TextOption: TextOption{
								Order:         int(qst["order"].(float64)),
								Content:       qst["content"].(string),
								Mark:          int(qst["mark"].(float64)),
								CaseSensitive: qst["case_sensitive"].(bool),
							},
						}, qRes.ID, qRes.QuestionHistoryID, userID)
						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
							return
						}
					} else if qRes.Type == util.Matching {
						if qst["type"].(string) != "MATCHING_ANSWER" {
							_, err := h.Service.CreateMatchingOption(c.Request.Context(), tx, &MatchingOptionRequest{
								MatchingOption: MatchingOption{
									Order:     int(qst["order"].(float64)),
									Content:   qst["content"].(string),
									Type:      qst["type"].(string),
									Eliminate: qst["eliminate"].(bool),
								},
							}, qRes.ID, qRes.QuestionHistoryID, userID)

							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

						} else {

							prompt, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c.Request.Context(), qRes.ID, int(qst["prompt"].(float64)))
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							option, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c.Request.Context(), qRes.ID, int(qst["option"].(float64)))
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							_, err = h.Service.CreateMatchingAnswer(c.Request.Context(), tx, &MatchingAnswerRequest{
								MatchingAnswer: MatchingAnswer{
									PromptID: prompt.ID,
									OptionID: option.ID,
									Mark:     int(qst["mark"].(float64)),
								},
							}, qRes.ID, qRes.QuestionHistoryID, userID)

							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}
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

	}

	h.Service.CommitTransaction(c.Request.Context(), tx)

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

func (h *Handler) DeleteQuiz(c *gin.Context) {
	_, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	quizID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	tx, err := h.Service.BeginTransaction(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questionPoolData, err := h.Service.GetQuestionPoolsByQuizID(c.Request.Context(), quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, questionPool := range questionPoolData {
		err := h.Service.DeleteQuestionPool(c.Request.Context(), tx, questionPool.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	questionData, err := h.Service.GetQuestionsByQuizID(c.Request.Context(), quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, question := range questionData {

		if question.Type == util.Choice || question.Type == util.TrueFalse {
			choiceOptionData, err := h.Service.GetChoiceOptionsByQuestionID(c.Request.Context(), question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			for _, choice := range choiceOptionData {
				err := h.Service.DeleteChoiceOption(c.Request.Context(), tx, choice.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

			}
		}

		if question.Type == util.ShortText || question.Type == util.Paragraph {
			textOptionData, err := h.Service.GetTextOptionsByQuestionID(c.Request.Context(), question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			for _, text := range textOptionData {
				err := h.Service.DeleteTextOption(c.Request.Context(), tx, text.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}

		if question.Type == util.Matching {
			matchingOptionData, err := h.Service.GetMatchingOptionsByQuestionID(c.Request.Context(), question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			matchinAnswerData, err := h.Service.GetMatchingAnswersByQuestionID(c.Request.Context(), question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			for _, matchingOption := range matchingOptionData {
				err := h.Service.DeleteMatchingOption(c.Request.Context(), tx, matchingOption.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}

			for _, matchingAnswer := range matchinAnswerData {
				err := h.Service.DeleteMatchingAnswer(c.Request.Context(), tx, matchingAnswer.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}

		err := h.Service.DeleteQuestion(c.Request.Context(), tx, question.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}

	err = h.Service.DeleteQuiz(c.Request.Context(), tx, quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Service.CommitTransaction(c.Request.Context(), tx)

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted",
	})
}

func (h *Handler) RestoreQuiz(c *gin.Context) {
	_, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	quizID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tx, err := h.Service.BeginTransaction(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questionPoolData, err := h.Service.GetDeleteQuestionPoolsByQuizID(c.Request.Context(), quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, questionPool := range questionPoolData {
		err := h.Service.RestoreQuestionPool(c.Request.Context(), tx, questionPool.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	questionData, err := h.Service.GetDeleteQuestionsByQuizID(c.Request.Context(), quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, question := range questionData {

		if question.Type == util.Choice || question.Type == util.TrueFalse {
			choiceOptionData, err := h.Service.GetDeleteChoiceOptionsByQuestionID(c.Request.Context(), question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			for _, choice := range choiceOptionData {
				err := h.Service.RestoreChoiceOption(c.Request.Context(), tx, choice.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

			}
		}

		if question.Type == util.ShortText || question.Type == util.Paragraph {
			textOptionData, err := h.Service.GetDeleteTextOptionsByQuestionID(c.Request.Context(), question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			for _, text := range textOptionData {
				err := h.Service.RestoreTextOption(c.Request.Context(), tx, text.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}

		if question.Type == util.Matching {
			matchingOptionData, err := h.Service.GetDeleteMatchingOptionsByQuestionID(c.Request.Context(), question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			matchinAnswerData, err := h.Service.GetDeleteMatchingAnswersByQuestionID(c.Request.Context(), question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			for _, matchingOption := range matchingOptionData {
				err := h.Service.RestoreMatchingOption(c.Request.Context(), tx, matchingOption.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}

			for _, matchingAnswer := range matchinAnswerData {
				err := h.Service.RestoreMatchingAnswer(c.Request.Context(), tx, matchingAnswer.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}

		err := h.Service.RestoreQuestion(c.Request.Context(), tx, question.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}

	err = h.Service.RestoreQuiz(c.Request.Context(), tx, quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Service.CommitTransaction(c.Request.Context(), tx)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Restore Quiz",
	})
}

// Quiz related with History Page
func (h *Handler) GetQuizHistories(c *gin.Context) {
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

	res, err := h.Service.GetQuizHistories(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	r := make([]QuizHistoryResponse, 0)

	for _, q := range res {
		qpRes, err := h.Service.GetQuestionPoolHistoriesByQuizID(c.Request.Context(), q.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Put All Question Pools Data in Questions Response
		for _, qpr := range qpRes {
			q.QuestionHistory = append(q.QuestionHistory, QuestionHistoryResponse{
				QuestionHistory: QuestionHistory{
					ID:             qpr.ID,
					QuestionID:     qpr.QuestionPoolID, //This Put QuestionPoolID in QuestionID instead
					QuizID:         qpr.QuizID,
					Type:           "POOL",
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

		qRes, err := h.Service.GetQuestionHistoriesByQuizID(c.Request.Context(), q.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, qr := range qRes {

			if qr.Type == util.Choice || qr.Type == util.TrueFalse {
				ocRes, err := h.Service.GetChoiceOptionHistoriesByQuestionID(c.Request.Context(), qr.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				var oc []any
				for _, ocr := range ocRes {
					oc = append(oc, ChoiceOptionHistoryResponse{
						ChoiceOptionHistory: ChoiceOptionHistory{
							ID:             ocr.ID,
							ChoiceOptionID: ocr.ChoiceOptionID,
							QuestionID:     ocr.QuestionID,
							Order:          ocr.Order,
							Content:        ocr.Content,
							Mark:           ocr.Mark,
							Color:          ocr.Color,
							Correct:        ocr.Correct,
							CreatedAt:      ocr.CreatedAt,
							UpdatedAt:      ocr.UpdatedAt,
							DeletedAt:      ocr.DeletedAt,
						},
					})
				}

				q.QuestionHistory = append(q.QuestionHistory, QuestionHistoryResponse{
					QuestionHistory: QuestionHistory{
						ID:             qr.ID,
						QuestionID:     qr.QuestionID,
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
				otRes, err := h.Service.GetTextOptionHistoriesByQuestionID(c.Request.Context(), qr.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				var ot []any
				for _, otr := range otRes {
					ot = append(ot, TextOptionHistoryResponse{
						TextOptionHistory: TextOptionHistory{
							ID:            otr.ID,
							OptionTextID:  otr.OptionTextID,
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

				q.QuestionHistory = append(q.QuestionHistory, QuestionHistoryResponse{
					QuestionHistory: QuestionHistory{
						ID:             qr.ID,
						QuestionID:     qr.QuestionID,
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
			} else if qr.Type == util.Matching {
				omRes, err := h.Service.GetMatchingOptionHistoriesByQuestionID(c.Request.Context(), qr.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				amRes, err := h.Service.GetMatchingAnswerHistoriesByQuestionID(c.Request.Context(), qr.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				var o []any
				for _, omr := range omRes {
					o = append(o, MatchingOptionAndAnswerHistoryResponse{
						ID:               omr.ID,
						OptionMatchingID: omr.OptionMatchingID,
						QuestionID:       omr.QuestionID,
						Type:             omr.Type,
						Order:            omr.Order,
						Content:          omr.Content,
						Eliminate:        omr.Eliminate,
						PromptID:         uuid.Nil,
						OptionID:         uuid.Nil,
						Mark:             0,
						CreatedAt:        omr.CreatedAt,
						UpdatedAt:        omr.UpdatedAt,
						DeletedAt:        omr.DeletedAt,
					})
				}

				for _, amr := range amRes {
					o = append(o, MatchingOptionAndAnswerResponse{
						ID:         amr.ID,
						QuestionID: amr.QuestionID,
						Type:       "MATCHING_ANSWER",
						Order:      0,
						Content:    "",
						Eliminate:  false,
						PromptID:   amr.PromptID,
						OptionID:   amr.OptionID,
						Mark:       amr.Mark,
						CreatedAt:  amr.CreatedAt,
						UpdatedAt:  amr.UpdatedAt,
						DeletedAt:  amr.DeletedAt,
					})
				}

				q.QuestionHistory = append(q.QuestionHistory, QuestionHistoryResponse{
					QuestionHistory: QuestionHistory{
						ID:             qr.ID,
						QuestionID:     qr.QuestionID,
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
					Options: o,
				})
			}
		}

		r = append(r, q)
	}

	// Bubble Sort the Questions and Question Pool by Order
	for _, res := range r {
		n := len(res.QuestionHistory)
		for i := 0; i < n-1; i++ {
			for j := 0; j < n-i-1; j++ {
				if res.QuestionHistory[j].Order > res.QuestionHistory[j+1].Order {
					res.QuestionHistory[j], res.QuestionHistory[j+1] = res.QuestionHistory[j+1], res.QuestionHistory[j]
				}
			}
		}
	}

	c.JSON(http.StatusOK, r)
}

func (h *Handler) GetQuizHistoryByID(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	res, err := h.Service.GetQuizHistoryByID(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	qpRes, err := h.Service.GetQuestionPoolHistoriesByQuizID(c.Request.Context(), res.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Put All Question Pools Data in Questions Response
	for _, qpr := range qpRes {
		res.QuestionHistory = append(res.QuestionHistory, QuestionHistoryResponse{
			QuestionHistory: QuestionHistory{
				ID:             qpr.ID,
				QuestionID:     qpr.QuestionPoolID,
				QuizID:         qpr.QuizID,
				Type:           "POOL",
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

	qRes, err := h.Service.GetQuestionHistoriesByQuizID(c.Request.Context(), res.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, qr := range qRes {
		if qr.Type == util.Choice || qr.Type == util.TrueFalse {
			ocRes, err := h.Service.GetChoiceOptionHistoriesByQuestionID(c.Request.Context(), qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var oc []any
			for _, ocr := range ocRes {
				oc = append(oc, ChoiceOptionHistoryResponse{
					ChoiceOptionHistory: ChoiceOptionHistory{
						ID:             ocr.ID,
						ChoiceOptionID: ocr.ChoiceOptionID,
						QuestionID:     ocr.QuestionID,
						Order:          ocr.Order,
						Content:        ocr.Content,
						Mark:           ocr.Mark,
						Color:          ocr.Color,
						Correct:        ocr.Correct,
						CreatedAt:      ocr.CreatedAt,
						UpdatedAt:      ocr.UpdatedAt,
						DeletedAt:      ocr.DeletedAt,
					},
				})
			}

			res.QuestionHistory = append(res.QuestionHistory, QuestionHistoryResponse{
				QuestionHistory: QuestionHistory{
					ID:             qr.ID,
					QuestionID:     qr.QuestionID,
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
			otRes, err := h.Service.GetTextOptionHistoriesByQuestionID(c.Request.Context(), qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var ot []any
			for _, otr := range otRes {
				ot = append(ot, TextOptionHistoryResponse{
					TextOptionHistory: TextOptionHistory{
						ID:            otr.ID,
						OptionTextID:  otr.OptionTextID,
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

			res.QuestionHistory = append(res.QuestionHistory, QuestionHistoryResponse{
				QuestionHistory: QuestionHistory{
					ID:             qr.ID,
					QuestionID:     qr.QuestionID,
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
		} else if qr.Type == util.Matching {
			omRes, err := h.Service.GetMatchingOptionHistoriesByQuestionID(c.Request.Context(), qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			amRes, err := h.Service.GetMatchingAnswerHistoriesByQuestionID(c.Request.Context(), qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var o []any
			for _, omr := range omRes {
				o = append(o, MatchingOptionAndAnswerHistoryResponse{
					ID:               omr.ID,
					OptionMatchingID: omr.OptionMatchingID,
					QuestionID:       omr.QuestionID,
					Type:             omr.Type,
					Order:            omr.Order,
					Content:          omr.Content,
					Eliminate:        omr.Eliminate,
					PromptID:         uuid.Nil,
					OptionID:         uuid.Nil,
					Mark:             0,
					CreatedAt:        omr.CreatedAt,
					UpdatedAt:        omr.UpdatedAt,
					DeletedAt:        omr.DeletedAt,
				})
			}

			for _, amr := range amRes {
				o = append(o, MatchingOptionAndAnswerHistoryResponse{
					ID:               amr.ID,
					AmswerMatchingID: amr.AnswerMatchingID,
					QuestionID:       amr.QuestionID,
					Type:             "MATCHING_ANSWER",
					Order:            0,
					Content:          "",
					Eliminate:        false,
					PromptID:         amr.PromptID,
					OptionID:         amr.OptionID,
					Mark:             amr.Mark,
					CreatedAt:        amr.CreatedAt,
					UpdatedAt:        amr.UpdatedAt,
					DeletedAt:        amr.DeletedAt,
				})
			}

			res.QuestionHistory = append(res.QuestionHistory, QuestionHistoryResponse{
				QuestionHistory: QuestionHistory{
					ID:             qr.ID,
					QuestionID:     qr.QuestionID,
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
				Options: o,
			})
		}
	}

	// Bubble Sort the Questions and Question Pool by Order
	n := len(res.QuestionHistory)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if res.QuestionHistory[j].Order > res.QuestionHistory[j+1].Order {
				res.QuestionHistory[j], res.QuestionHistory[j+1] = res.QuestionHistory[j+1], res.QuestionHistory[j]
			}
		}
	}

	c.JSON(http.StatusOK, res)
}
