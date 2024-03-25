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
	tx, _ := h.Service.BeginTransaction(c)

	res, err := h.Service.CreateQuiz(c, tx, &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Service.CommitTransaction(c, tx)

	var qpResID *uuid.UUID = nil
	var qphResID *uuid.UUID = nil

	for _, q := range req.Questions {
		var qRes *CreateQuestionResponse
		var qpRes *CreateQuestionPoolResponse

		if q.Type == util.Pool {
			txPool, _ := h.Service.BeginTransaction(c)

			qpRes, err = h.Service.CreateQuestionPool(c, txPool, &q, res.ID, res.QuizHistoryID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			h.Service.CommitTransaction(c, txPool)
			qpResID = &qpRes.ID
			qphResID = &qpRes.QuestionPoolHistoryID

			res.Questions = append(res.Questions, QuestionResponse{
				Question: Question{
					ID:        qpRes.ID,
					QuizID:    qpRes.QuizID,
					Type:      "POOL",
					Order:     qpRes.Order,
					PoolOrder: qpRes.PoolOrder,
					// PoolRequired:   false,
					Content:        qpRes.Content,
					Note:           qpRes.Note,
					Media:          qpRes.Media,
					MediaType:      qpRes.MediaType,
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
			txQuestion, _ := h.Service.BeginTransaction(c)
			if q.IsInPool {
				qRes, err = h.Service.CreateQuestion(c, txQuestion, &q, res.ID, res.QuizHistoryID, qpResID, qphResID, userID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			} else {
				qRes, err = h.Service.CreateQuestion(c, txQuestion, &q, res.ID, res.QuizHistoryID, nil, nil, userID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
			h.Service.CommitTransaction(c, txQuestion)

			for _, qt := range qRes.Options {
				if qst, ok := qt.(map[string]any); ok {
					if qRes.Type == util.Choice || qRes.Type == util.TrueFalse {
						txChoice, _ := h.Service.BeginTransaction(c)
						_, err := h.Service.CreateChoiceOption(c, txChoice, &ChoiceOptionRequest{
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
						h.Service.CommitTransaction(c, txChoice)

					} else if qRes.Type == util.FillBlank || qRes.Type == util.Paragraph {
						txText, _ := h.Service.BeginTransaction(c)
						_, err := h.Service.CreateTextOption(c, txText, &TextOptionRequest{
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

						h.Service.CommitTransaction(c, txText)

					} else if qRes.Type == util.Matching {
						txMatching, _ := h.Service.BeginTransaction(c)
						if qst["type"].(string) != "MATCHING_ANSWER" {
							_, err := h.Service.CreateMatchingOption(c, txMatching, &MatchingOptionRequest{
								MatchingOption: MatchingOption{
									Order:     int(qst["order"].(float64)),
									Content:   qst["content"].(string),
									Type:      qst["type"].(string),
									Color:     qst["color"].(string),
									Eliminate: qst["eliminate"].(bool),
								},
							}, qRes.ID, qRes.QuestionHistoryID, userID)

							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

						} else {

							prompt, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c, qRes.ID, int(qst["prompt"].(float64)))
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							option, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c, qRes.ID, int(qst["option"].(float64)))
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							promptH, err := h.Service.GetMatchingOptionHistoryByOptionMatchingID(c, prompt.ID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							optionH, err := h.Service.GetMatchingOptionHistoryByOptionMatchingID(c, option.ID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							_, err = h.Service.CreateMatchingAnswer(c, txMatching, &MatchingAnswerRequest{
								MatchingAnswer: MatchingAnswer{
									PromptID: prompt.ID,
									OptionID: option.ID,
									Mark:     int(qst["mark"].(float64)),
								},
							}, qRes.ID, qRes.QuestionHistoryID, promptH.ID, optionH.ID, userID)

							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}
						}
						h.Service.CommitTransaction(c, txMatching)

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
				PoolOrder:      qRes.PoolOrder,
				PoolRequired:   qRes.PoolRequired,
				Content:        qRes.Content,
				Note:           qRes.Note,
				Media:          qRes.Media,
				MediaType:      qRes.MediaType,
				UseTemplate:    qRes.UseTemplate,
				TimeLimit:      qRes.TimeLimit,
				HaveTimeFactor: qRes.HaveTimeFactor,
				TimeFactor:     qRes.TimeFactor,
				FontSize:       qRes.FontSize,
				LayoutIdx:      qRes.LayoutIdx,
				SelectMin:      qRes.SelectMin,
				SelectMax:      qRes.SelectMax,
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
			SelectMin:      res.SelectMin,
			SelectMax:      res.SelectMax,
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

	res, err := h.Service.GetQuizzes(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
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

	res, err := h.Service.GetQuizByID(c, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	tx, err := h.Service.BeginTransaction(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.UpdateQuiz(c, tx, &req, userID, quizID)
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
				qpRes, err = h.Service.UpdateQuestionPool(c, tx, &q, userID, q.ID, res.QuizHistoryID)
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
						ID:        qpRes.ID,
						QuizID:    qpRes.QuizID,
						Type:      "POOL",
						Order:     qpRes.Order,
						PoolOrder: qpRes.PoolOrder,
						// PoolRequired:   qpRes.PoolRequired,
						Content:        qpRes.Content,
						Note:           qpRes.Note,
						Media:          qpRes.Media,
						MediaType:      qpRes.MediaType,
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
				if q.IsInPool {
					qRes, err = h.Service.UpdateQuestion(c, tx, &q, userID, q.ID, res.QuizHistoryID, qphResID)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
				} else {
					qRes, err = h.Service.UpdateQuestion(c, tx, &q, userID, q.ID, res.QuizHistoryID, nil)
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
							_, err := h.Service.UpdateChoiceOption(c, tx, &choiceReq, userID, id, qRes.QuestionHistoryID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

						} else {
							_, err := h.Service.CreateChoiceOption(c, tx, &choiceReq, qRes.ID, qRes.QuestionHistoryID, userID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}
						}

					} else if qRes.Type == util.FillBlank || qRes.Type == util.Paragraph {
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
							_, err := h.Service.UpdateTextOption(c, tx, &textReq, userID, id, qRes.QuestionHistoryID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

						} else {
							_, err := h.Service.CreateTextOption(c, tx, &textReq, qRes.ID, qRes.QuestionHistoryID, userID)
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
									Color:      qst["color"].(string),
									Eliminate:  qst["eliminate"].(bool),
								},
							}

							if matchingOptionReq.ID != uuid.Nil {
								_, err := h.Service.UpdateMatchingOption(c, tx, &matchingOptionReq, userID, id, qRes.QuestionHistoryID)
								if err != nil {
									c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
									return
								}

							} else {
								_, err := h.Service.CreateMatchingOption(c, tx, &matchingOptionReq, qRes.ID, qRes.QuestionHistoryID, userID)
								if err != nil {
									c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
									return
								}
							}
						} else {
							prompt, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c, qRes.ID, int(qst["prompt"].(float64)))
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							option, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c, qRes.ID, int(qst["option"].(float64)))
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							promptH, err := h.Service.GetMatchingOptionHistoryByOptionMatchingID(c, prompt.ID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							optionH, err := h.Service.GetMatchingOptionHistoryByOptionMatchingID(c, option.ID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							if id != uuid.Nil {
								_, err = h.Service.UpdateMatchingAnswer(c, tx, &MatchingAnswerRequest{
									MatchingAnswer: MatchingAnswer{
										ID:         id,
										QuestionID: questionID,
										PromptID:   prompt.ID,
										OptionID:   option.ID,
										Mark:       int(qst["mark"].(float64)),
									},
								}, userID, id, qRes.QuestionHistoryID)

								if err != nil {
									c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
									return
								}
							} else {
								_, err = h.Service.CreateMatchingAnswer(c, tx, &MatchingAnswerRequest{
									MatchingAnswer: MatchingAnswer{
										PromptID: prompt.ID,
										OptionID: option.ID,
										Mark:     int(qst["mark"].(float64)),
									},
								}, qRes.ID, qRes.QuestionHistoryID, promptH.ID, optionH.ID, userID)

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
					PoolOrder:      qRes.PoolOrder,
					PoolRequired:   qRes.PoolRequired,
					Content:        qRes.Content,
					Note:           qRes.Note,
					Media:          qRes.Media,
					MediaType:      qRes.MediaType,
					UseTemplate:    qRes.UseTemplate,
					TimeLimit:      qRes.TimeLimit,
					HaveTimeFactor: qRes.HaveTimeFactor,
					TimeFactor:     qRes.TimeFactor,
					FontSize:       qRes.FontSize,
					LayoutIdx:      qRes.LayoutIdx,
					SelectMin:      qRes.SelectMin,
					SelectMax:      qRes.SelectMax,
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
				qpRes, err = h.Service.CreateQuestionPool(c, tx, &q, quizID, res.QuizHistoryID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				qpResID = &qpRes.ID
				qphResID = &qpRes.QuestionPoolHistoryID
				continue
			}

			if q.IsInPool {
				qRes, err = h.Service.CreateQuestion(c, tx, &q, quizID, res.QuizHistoryID, qpResID, qphResID, userID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			} else {
				qRes, err = h.Service.CreateQuestion(c, tx, &q, quizID, res.QuizHistoryID, nil, nil, userID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}

			for _, qt := range qRes.Options {
				if qst, ok := qt.(map[string]any); ok {
					if qRes.Type == util.Choice || qRes.Type == util.TrueFalse {
						_, err := h.Service.CreateChoiceOption(c, tx, &ChoiceOptionRequest{
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
					} else if qRes.Type == util.FillBlank || qRes.Type == util.Paragraph {
						_, err := h.Service.CreateTextOption(c, tx, &TextOptionRequest{
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
							_, err := h.Service.CreateMatchingOption(c, tx, &MatchingOptionRequest{
								MatchingOption: MatchingOption{
									Order:     int(qst["order"].(float64)),
									Content:   qst["content"].(string),
									Type:      qst["type"].(string),
									Color:     qst["color"].(string),
									Eliminate: qst["eliminate"].(bool),
								},
							}, qRes.ID, qRes.QuestionHistoryID, userID)

							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

						} else {

							prompt, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c, qRes.ID, int(qst["prompt"].(float64)))
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							option, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c, qRes.ID, int(qst["option"].(float64)))
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							promptH, err := h.Service.GetMatchingOptionHistoryByOptionMatchingID(c, prompt.ID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							optionH, err := h.Service.GetMatchingOptionHistoryByOptionMatchingID(c, option.ID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							_, err = h.Service.CreateMatchingAnswer(c, tx, &MatchingAnswerRequest{
								MatchingAnswer: MatchingAnswer{
									PromptID: prompt.ID,
									OptionID: option.ID,
									Mark:     int(qst["mark"].(float64)),
								},
							}, qRes.ID, qRes.QuestionHistoryID, promptH.ID, optionH.ID, userID)

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
					PoolOrder:      qRes.PoolOrder,
					PoolRequired:   qRes.PoolRequired,
					Content:        qRes.Content,
					Note:           qRes.Note,
					Media:          qRes.Media,
					MediaType:      qRes.MediaType,
					UseTemplate:    qRes.UseTemplate,
					TimeLimit:      qRes.TimeLimit,
					HaveTimeFactor: qRes.HaveTimeFactor,
					TimeFactor:     qRes.TimeFactor,
					FontSize:       qRes.FontSize,
					LayoutIdx:      qRes.LayoutIdx,
					SelectMin:      qRes.SelectMin,
					SelectMax:      qRes.SelectMax,
					CreatedAt:      qRes.CreatedAt,
					UpdatedAt:      qRes.UpdatedAt,
					DeletedAt:      qRes.DeletedAt,
				},
				Options: qRes.Options,
			})

		}

	}

	h.Service.CommitTransaction(c, tx)

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
			SelectMin:      res.SelectMin,
			SelectMax:      res.SelectMax,
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

	tx, err := h.Service.BeginTransaction(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questionPoolData, err := h.Service.GetQuestionPoolsByQuizID(c, quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, questionPool := range questionPoolData {
		err := h.Service.DeleteQuestionPool(c, tx, questionPool.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	questionData, err := h.Service.GetQuestionsByQuizID(c, quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, question := range questionData {

		if question.Type == util.Choice || question.Type == util.TrueFalse {
			choiceOptionData, err := h.Service.GetChoiceOptionsByQuestionID(c, question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			for _, choice := range choiceOptionData {
				err := h.Service.DeleteChoiceOption(c, tx, choice.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

			}
		}

		if question.Type == util.FillBlank || question.Type == util.Paragraph {
			textOptionData, err := h.Service.GetTextOptionsByQuestionID(c, question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			for _, text := range textOptionData {
				err := h.Service.DeleteTextOption(c, tx, text.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}

		if question.Type == util.Matching {
			matchingOptionData, err := h.Service.GetMatchingOptionsByQuestionID(c, question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			matchinAnswerData, err := h.Service.GetMatchingAnswersByQuestionID(c, question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			for _, matchingOption := range matchingOptionData {
				err := h.Service.DeleteMatchingOption(c, tx, matchingOption.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}

			for _, matchingAnswer := range matchinAnswerData {
				err := h.Service.DeleteMatchingAnswer(c, tx, matchingAnswer.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}

		err := h.Service.DeleteQuestion(c, tx, question.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}

	err = h.Service.DeleteQuiz(c, tx, quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Service.CommitTransaction(c, tx)

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

	tx, err := h.Service.BeginTransaction(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questionPoolData, err := h.Service.GetDeleteQuestionPoolsByQuizID(c, quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, questionPool := range questionPoolData {
		err := h.Service.RestoreQuestionPool(c, tx, questionPool.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	questionData, err := h.Service.GetDeleteQuestionsByQuizID(c, quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, question := range questionData {

		if question.Type == util.Choice || question.Type == util.TrueFalse {
			choiceOptionData, err := h.Service.GetDeleteChoiceOptionsByQuestionID(c, question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			for _, choice := range choiceOptionData {
				err := h.Service.RestoreChoiceOption(c, tx, choice.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

			}
		}

		if question.Type == util.FillBlank || question.Type == util.Paragraph {
			textOptionData, err := h.Service.GetDeleteTextOptionsByQuestionID(c, question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			for _, text := range textOptionData {
				err := h.Service.RestoreTextOption(c, tx, text.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}

		if question.Type == util.Matching {
			matchingOptionData, err := h.Service.GetDeleteMatchingOptionsByQuestionID(c, question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			matchinAnswerData, err := h.Service.GetDeleteMatchingAnswersByQuestionID(c, question.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			for _, matchingOption := range matchingOptionData {
				err := h.Service.RestoreMatchingOption(c, tx, matchingOption.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}

			for _, matchingAnswer := range matchinAnswerData {
				err := h.Service.RestoreMatchingAnswer(c, tx, matchingAnswer.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}

		err := h.Service.RestoreQuestion(c, tx, question.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}

	err = h.Service.RestoreQuiz(c, tx, quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Service.CommitTransaction(c, tx)

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

	res, err := h.Service.GetQuizHistories(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
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

	res, err := h.Service.GetQuizHistoryByID(c, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
