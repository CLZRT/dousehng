package service

import (
	"demo1/model"
	"demo1/repository"
	"errors"
	"fmt"
	"time"
)

// AddComment 登录用户对视频进行评论
func AddComment(req *model.CommentActionRequest) (*model.CommentActionResponse, error) {
	// 创建单例
	commentDao := repository.NewCommentDAO()
	userDao := repository.NewUserDAO()

	commentId, err := commentDao.CreateComment(req.UserID, req.VideoID, &req.CommentText)
	if err != nil {
		return &model.CommentActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "create comment error",
			},
			Comment: model.Comment{},
		}, errors.New("create comment error")
	}

	var author model.User
	if err := userDao.FindUserById(req.UserID, (*repository.User)(&author)); err != nil { // 找作者信息
		return &model.CommentActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "comment author not exists",
			},
			Comment: model.Comment{},
		}, errors.New("comment author not exists")
	}

	return &model.CommentActionResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "ok",
		},
		Comment: model.Comment{
			ID:         commentId,
			User:       author,
			Content:    req.CommentText,
			CreateDate: time.Now().Format("01-02"),
		},
	}, nil
}

func DeleteComment(req *model.CommentActionRequest) (*model.CommentActionResponse, error) {
	// 创建单例
	commentDao := repository.NewCommentDAO()

	err := commentDao.DeleteCommentById(req.CommentID)
	if err != nil {
		return &model.CommentActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "delete comment error",
			},
			Comment: model.Comment{},
		}, errors.New("delete comment error")
	}

	return &model.CommentActionResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "ok",
		},
		Comment: model.Comment{},
	}, errors.New("comment delete ok")
}

// CommentList 查看视频的所有评论，按发布时间倒序
func CommentList(req *model.CommentListRequest) (*model.CommentListResponse, error) {

	// 创建单例
	commentDAO := repository.NewCommentDAO()
	userDao := repository.NewUserDAO()

	var commentList []repository.Comment // 猜想，如果评论量特别大的话，是不是可以做成分段查询的，how，是不是需要前端来请求
	if err := commentDAO.GetAllComment(&commentList, req.VideoID); err != nil {
		return &model.CommentListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "get comment list error",
			},
			CommentList: nil,
		}, errors.New("get comment list error")
	}

	fmt.Println("comment list length : ", len(commentList))
	//for _, comment := range commentList {
	//	fmt.Printf("%+v\n", comment)
	//}

	// 构造结果
	var author model.User
	resList := make([]model.Comment, len(commentList))
	for i, comment := range commentList {
		// 找到评论作者信息
		if err := userDao.FindUserById(comment.AuthorID, (*repository.User)(&author)); err != nil {
			return &model.CommentListResponse{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  "comment author not exists",
				},
				CommentList: nil,
			}, errors.New("comment author not exists")
		}
		resList[i] = model.Comment{
			ID:         comment.ID,
			User:       author, // 如果这个能直接成功拿到的话，前面有一个通过id来找人的逻辑就不用写了
			Content:    comment.Content,
			CreateDate: time.Unix(comment.CommentPublishTime, 0).Format("01-02"),
		}
	}

	for _, comment := range resList {
		fmt.Printf("%+v\n", comment)
	}

	// 返回结果
	return &model.CommentListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "ok",
		},
		CommentList: &resList,
	}, nil
}
