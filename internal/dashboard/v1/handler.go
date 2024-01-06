package v1

import (
	"net/http"

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

// func (h *Handler) CreateAnswerResponse(c *gin.Context) {

// 	var req LiveAnswerRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
	
// 	uid, ok := c.Get("uid")
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}

// 	userID, err := uuid.Parse(uid.(string))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
// 		return
// 	}

// 	res, er := h.Service.CreateAnswerResponse(c.Request.Context(), &req)
// 	if er != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": er.Error()})
// 		return
// 	}

// 	for _, answer := range {

// 	}

	

// 	for _, q := range req.Questions {
// 		var qRes *CreateQuestionResponse
// 		var qpRes *CreateQuestionPoolResponse

// 		if q.Type == util.Pool {
// 			qpRes, err = h.Service.CreateQuestionPool(c.Request.Context(), &q, res.ID, res.QuizHistoryID)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 				return
// 			}
// 			qpResID = &qpRes.ID
// 			qphResID = &qpRes.QuestionPoolHistoryID

// 			res.Questions = append(res.Questions, QuestionResponse{
// 				Question: Question{
// 					ID:             qpRes.ID,
// 					QuizID:         qpRes.QuizID,
// 					Type:           "POOL",
// 					Order:          qpRes.Order,
// 					Content:        qpRes.Content,
// 					Note:           qpRes.Note,
// 					Media:          qpRes.Media,
// 					TimeLimit:      qpRes.TimeLimit,
// 					HaveTimeFactor: qpRes.HaveTimeFactor,
// 					TimeFactor:     qpRes.TimeFactor,
// 					FontSize:       qpRes.FontSize,
// 					CreatedAt:      qpRes.CreatedAt,
// 					UpdatedAt:      qpRes.UpdatedAt,
// 					DeletedAt:      qpRes.DeletedAt,
// 				},
// 			})

// 			continue
// 		} else {
// 			if q.IsInPool == true {
// 				qRes, err = h.Service.CreateQuestion(c.Request.Context(), &q, res.ID, res.QuizHistoryID, qpResID, qphResID, userID)
// 				if err != nil {
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 					return
// 				}
// 			} else {
// 				qRes, err = h.Service.CreateQuestion(c.Request.Context(), &q, res.ID, res.QuizHistoryID, nil, nil, userID)
// 				if err != nil {
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 					return
// 				}
// 			}

// 			for _, qt := range qRes.Options {
// 				if qst, ok := qt.(map[string]any); ok {
// 					if qRes.Type == util.Choice || qRes.Type == util.TrueFalse {
// 						_, err := h.Service.CreateChoiceOption(c.Request.Context(), &ChoiceOptionRequest{
// 							ChoiceOption: ChoiceOption{
// 								Order:   int(qst["order"].(float64)),
// 								Content: qst["content"].(string),
// 								Mark:    int(qst["mark"].(float64)),
// 								Color:   qst["color"].(string),
// 								Correct: qst["correct"].(bool),
// 							},
// 						}, qRes.ID, qRes.QuestionHistoryID, userID)

// 						if err != nil {
// 							c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 							return
// 						}

// 					} else if qRes.Type == util.ShortText || qRes.Type == util.Paragraph {
// 						_, err := h.Service.CreateTextOption(c.Request.Context(), &TextOptionRequest{
// 							TextOption: TextOption{
// 								Order:         int(qst["order"].(float64)),
// 								Content:       qst["content"].(string),
// 								Mark:          int(qst["mark"].(float64)),
// 								CaseSensitive: qst["case_sensitive"].(bool),
// 							},
// 						}, qRes.ID, qRes.QuestionHistoryID, userID)

// 						if err != nil {
// 							c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 							return
// 						}

// 					} else if qRes.Type == util.Matching {
// 						if qst["type"].(string) != "MATCHING_ANSWER" {
// 							_, err := h.Service.CreateMatchingOption(c.Request.Context(), &MatchingOptionRequest{
// 								MatchingOption: MatchingOption{
// 									Order:     int(qst["order"].(float64)),
// 									Content:   qst["content"].(string),
// 									Type:      qst["type"].(string),
// 									Eliminate: qst["eliminate"].(bool),
// 								},
// 							}, qRes.ID, qRes.QuestionHistoryID, userID)

// 							if err != nil {
// 								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 								return
// 							}

// 						} else {

// 							prompt, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c.Request.Context(), qRes.ID, int(qst["prompt"].(float64)))
// 							if err != nil {
// 								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 								return
// 							}

// 							option, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c.Request.Context(), qRes.ID, int(qst["option"].(float64)))
// 							if err != nil {
// 								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 								return
// 							}

// 							_, err = h.Service.CreateMatchingAnswer(c.Request.Context(), &MatchingAnswerRequest{
// 								MatchingAnswer: MatchingAnswer{
// 									PromptID: prompt.ID,
// 									OptionID: option.ID,
// 									Mark:     int(qst["mark"].(float64)),
// 								},
// 							}, qRes.ID, qRes.QuestionHistoryID, userID)

// 							if err != nil {
// 								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 								return
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}

// 		res.Questions = append(res.Questions, QuestionResponse{
// 			Question: Question{
// 				ID:             qRes.ID,
// 				QuizID:         qRes.QuizID,
// 				QuestionPoolID: qRes.QuestionPoolID,
// 				Type:           qRes.Type,
// 				Order:          qRes.Order,
// 				Content:        qRes.Content,
// 				Note:           qRes.Note,
// 				Media:          qRes.Media,
// 				UseTemplate:    qRes.UseTemplate,
// 				TimeLimit:      qRes.TimeLimit,
// 				HaveTimeFactor: qRes.HaveTimeFactor,
// 				TimeFactor:     qRes.TimeFactor,
// 				FontSize:       qRes.FontSize,
// 				LayoutIdx:      qRes.LayoutIdx,
// 				SelectUpTo:     qRes.SelectUpTo,
// 				CreatedAt:      qRes.CreatedAt,
// 				UpdatedAt:      qRes.UpdatedAt,
// 				DeletedAt:      qRes.DeletedAt,
// 			},
// 			Options: qRes.Options,
// 		})
// 	}

// 	c.JSON(http.StatusCreated, &QuizResponse{
// 		Quiz: Quiz{
// 			ID:             res.ID,
// 			CreatorID:      res.CreatorID,
// 			Title:          res.Title,
// 			Description:    res.Description,
// 			CoverImage:     res.CoverImage,
// 			Visibility:     res.Visibility,
// 			TimeLimit:      res.TimeLimit,
// 			HaveTimeFactor: res.HaveTimeFactor,
// 			TimeFactor:     res.TimeFactor,
// 			FontSize:       res.FontSize,
// 			Mark:           res.Mark,
// 			SelectUpTo:     res.SelectUpTo,
// 			CaseSensitive:  res.CaseSensitive,
// 			CreatedAt:      res.CreatedAt,
// 			UpdatedAt:      res.UpdatedAt,
// 			DeletedAt:      res.DeletedAt,
// 		},
// 		Questions: res.Questions,
// 	})
// }

func (h *Handler) GetAnswerResponseByLiveQuizSessionID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	_, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	res, err := h.Service.GetAnswerResponseByLiveQuizSessionID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}

func (h *Handler) GetAnswerResponseByQuestionID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	_, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	res, err := h.Service.GetAnswerResponseByQuestionID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}

func (h *Handler) GetAnswerResponseByParticipantID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	_, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	res, err := h.Service.GetAnswerResponseByParticipantID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}