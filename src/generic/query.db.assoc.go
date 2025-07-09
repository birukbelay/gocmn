package generic

import (
	"context"
	"errors"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/birukbelay/gocmn/src/dtos"
	respC "github.com/birukbelay/gocmn/src/resp_codes"
)

func CreateOneWithAssociations[T any, A any](u *gorm.DB, ctx context.Context, createDto, model2UpdateDto any, association AssocVar) (dtos.GResp[T], error) {
	tx := u.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return dtos.InternalErrS[T](err.Error()), err
	}
	targetModel := new(T)
	if err := mapstructure.Decode(createDto, &targetModel); err != nil {
		tx.Rollback()
		return dtos.BadReqErrS[T](err.Error()), err
	}

	query := tx.WithContext(ctx)
	if association.Debug {
		query = query.Debug()
	}
	for _, pre := range association.Preloads {
		query = query.Preload(pre)
	}
	resp := query.Clauses(clause.Returning{}).Create(&targetModel)
	if resp.Error != nil {
		tx.Rollback()
		return dtos.InternalErrS[T](resp.Error.Error()), resp.Error
	}

	//TODO how to fix it if the user wants it to be empty
	//using nil instead of len(arr) seemed to work until now
	if association.AssociatedValues != nil {
		query := "id IN (?)"
		if association.Key == AssociationName {
			query = "name IN (?)"
		} else if association.Key == AssociationId {
			query = "id IN (?)"
		} else if association.Key == AssociationSlug {
			query = "slug IN (?)"
		} else {
			tx.Rollback()
			return dtos.BadReqErrS[T](respC.DataNotFound.Msg()), errors.New(respC.DataNotFound.Msg())
		}
		var model2 A
		if err := tx.Model(&model2).Where(query, association.AssociatedValues).Updates(model2UpdateDto).Error; err != nil {
			tx.Rollback()
			return dtos.BadReqErrS[T](respC.DataNotFound.Msg()), errors.New(respC.DataNotFound.Msg())
		}
	}
	commit := tx.Commit()
	if commit.Error != nil {
		return dtos.InternalErrS[T](commit.Error.Error()), commit.Error
	}
	return dtos.SuccessCS[T](*targetModel, respC.UpdateSuccess, resp.RowsAffected), nil
}

func UpdateOneWithAssociations[T any, A any](u *gorm.DB, ctx context.Context, filter any, updateDto any, association AssocVar) (response dtos.GResp[T], err error) {
	//===========  Step1: define transaction
	tx := u.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			response, err = dtos.InternalErrS[T]("Transaction failed"), errors.New("Transaction Failed")
		}
	}()
	if err := tx.Error; err != nil {
		return dtos.InternalErrS[T](err.Error()), err
	}
	//=============== Step2: define query for first model
	result := new(T)
	query := tx.WithContext(ctx)
	if association.Debug {
		query = query.Debug()
	}
	for _, pre := range association.Preloads {
		query = query.Preload(pre)
	}
	resp := query.Clauses(clause.Returning{}).Model(&result).Where(filter).Updates(updateDto)
	if resp.Error != nil {
		tx.Rollback()
		return dtos.InternalErrS[T](resp.Error.Error()), resp.Error
	}
	if resp.RowsAffected == 0 {
		tx.Rollback()
		return dtos.NotFoundErrS[T](respC.NoRecordsUpdated.Msg()), errors.New(respC.NoRecordsUpdated.Msg())
	}

	//TODO how to fix it if the user wants it to be empty
	//using nil instead of len(arr) seemed to work until now
	if association.AssociatedValues != nil { //We use nill so that if user wants to remove all associations he can put nil

		secModQury := "id IN (?)"
		if association.Key == AssociationName {
			secModQury = "name IN (?)"
		} else if association.Key == AssociationId {
			secModQury = "id IN (?)"
		} else if association.Key == AssociationSlug {
			secModQury = "slug IN (?)"
		} else {
			tx.Rollback()
			return dtos.BadReqErrS[T](respC.DataNotFound.Msg()), errors.New(respC.DataNotFound.Msg())
		}
		//Query the associated model(eg: cards of users with `ID IN (cardIds)`)
		var associatedModels []*A
		if err := tx.Where(secModQury, association.AssociatedValues).Find(&associatedModels).Error; err != nil {
			tx.Rollback()
			return dtos.NotFoundErrS[T](respC.DataNotFound.Msg()), errors.New(respC.DataNotFound.Msg())
		}
		err := tx.Debug().Model(&result).Association(association.ModelName).Replace(associatedModels)
		if err != nil {
			tx.Rollback()
			return dtos.InternalErrS[T](respC.UpdatingAssociationsFailed.Msg()), errors.New(respC.UpdatingAssociationsFailed.Msg())
		}

	}
	commit := tx.Commit()
	if commit.Error != nil {
		return dtos.InternalErrS[T](commit.Error.Error()), commit.Error
	}
	return dtos.SuccessCS[T](*result, respC.UpdateSuccess, resp.RowsAffected), nil
}
