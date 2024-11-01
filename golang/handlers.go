package main

import (
	"encoding/json"
)

func GetProductComplex(req Req) {
	req.Respond(ProductComplexGet(req.IntParam))
}

func GetProductDeep(req Req) {
	counter := 1
	req.Respond(ProductDeepGet(req.IntParam, &counter))
}

func CreateProductToOrderingDefault(req Req) {
	p, err := DecodeProductToOrdering(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProductToOrderingCreateDefault(p))
}

func GetDocument(req Req) {
	req.Respond(DocumentGet(req.IntParam, nil))
}

func GetDocumentAll(req Req) {
	req.Respond(DocumentGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateDocument(req Req) {
	d, err := DecodeDocument(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(DocumentCreate(d, nil))
}

func UpdateDocument(req Req) {
	d, err := DecodeDocument(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(DocumentUpdate(d, nil))
}

func DeleteDocument(req Req) {
	req.Respond(DocumentDelete(req.IntParam, nil, false))
}

func GetDocumentByFilterInt(req Req) {
	req.Respond(DocumentGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetDocumentByFilterStr(req Req) {
	req.Respond(DocumentGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeDocument(req Req) (Document, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var d Document
	err := decoder.Decode(&d)
	return d, err
}

func GetMeasure(req Req) {
	req.Respond(MeasureGet(req.IntParam, nil))
}

func GetMeasureAll(req Req) {
	req.Respond(MeasureGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateMeasure(req Req) {
	m, err := DecodeMeasure(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MeasureCreate(m, nil))
}

func UpdateMeasure(req Req) {
	m, err := DecodeMeasure(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MeasureUpdate(m, nil))
}

func DeleteMeasure(req Req) {
	req.Respond(MeasureDelete(req.IntParam, nil, false))
}

func GetMeasureByFilterInt(req Req) {
	req.Respond(MeasureGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetMeasureByFilterStr(req Req) {
	req.Respond(MeasureGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeMeasure(req Req) (Measure, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var m Measure
	err := decoder.Decode(&m)
	return m, err
}

func GetCountType(req Req) {
	req.Respond(CountTypeGet(req.IntParam, nil))
}

func GetCountTypeAll(req Req) {
	req.Respond(CountTypeGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateCountType(req Req) {
	c, err := DecodeCountType(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(CountTypeCreate(c, nil))
}

func UpdateCountType(req Req) {
	c, err := DecodeCountType(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(CountTypeUpdate(c, nil))
}

func DeleteCountType(req Req) {
	req.Respond(CountTypeDelete(req.IntParam, nil, false))
}

func GetCountTypeByFilterInt(req Req) {
	req.Respond(CountTypeGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetCountTypeByFilterStr(req Req) {
	req.Respond(CountTypeGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeCountType(req Req) (CountType, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var c CountType
	err := decoder.Decode(&c)
	return c, err
}

func GetColorGroup(req Req) {
	req.Respond(ColorGroupGet(req.IntParam, nil))
}

func GetColorGroupAll(req Req) {
	req.Respond(ColorGroupGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateColorGroup(req Req) {
	c, err := DecodeColorGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ColorGroupCreate(c, nil))
}

func UpdateColorGroup(req Req) {
	c, err := DecodeColorGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ColorGroupUpdate(c, nil))
}

func DeleteColorGroup(req Req) {
	req.Respond(ColorGroupDelete(req.IntParam, nil, false))
}

func GetColorGroupByFilterInt(req Req) {
	req.Respond(ColorGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetColorGroupByFilterStr(req Req) {
	req.Respond(ColorGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeColorGroup(req Req) (ColorGroup, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var c ColorGroup
	err := decoder.Decode(&c)
	return c, err
}

func GetColor(req Req) {
	req.Respond(ColorGet(req.IntParam, nil))
}

func GetColorAll(req Req) {
	req.Respond(ColorGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateColor(req Req) {
	c, err := DecodeColor(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ColorCreate(c, nil))
}

func UpdateColor(req Req) {
	c, err := DecodeColor(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ColorUpdate(c, nil))
}

func DeleteColor(req Req) {
	req.Respond(ColorDelete(req.IntParam, nil, false))
}

func GetColorByFilterInt(req Req) {
	req.Respond(ColorGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetColorByFilterStr(req Req) {
	req.Respond(ColorGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeColor(req Req) (Color, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var c Color
	err := decoder.Decode(&c)
	return c, err
}

func GetMatherialGroup(req Req) {
	req.Respond(MatherialGroupGet(req.IntParam, nil))
}

func GetMatherialGroupAll(req Req) {
	req.Respond(MatherialGroupGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateMatherialGroup(req Req) {
	m, err := DecodeMatherialGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialGroupCreate(m, nil))
}

func UpdateMatherialGroup(req Req) {
	m, err := DecodeMatherialGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialGroupUpdate(m, nil))
}

func DeleteMatherialGroup(req Req) {
	req.Respond(MatherialGroupDelete(req.IntParam, nil, false))
}

func GetMatherialGroupByFilterInt(req Req) {
	req.Respond(MatherialGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetMatherialGroupByFilterStr(req Req) {
	req.Respond(MatherialGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeMatherialGroup(req Req) (MatherialGroup, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var m MatherialGroup
	err := decoder.Decode(&m)
	return m, err
}

func GetMatherial(req Req) {
	req.Respond(MatherialGet(req.IntParam, nil))
}

func GetMatherialAll(req Req) {
	req.Respond(MatherialGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateMatherial(req Req) {
	m, err := DecodeMatherial(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialCreate(m, nil))
}

func UpdateMatherial(req Req) {
	m, err := DecodeMatherial(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialUpdate(m, nil))
}

func DeleteMatherial(req Req) {
	req.Respond(MatherialDelete(req.IntParam, nil, false))
}

func GetMatherialByFilterInt(req Req) {
	req.Respond(MatherialGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetMatherialByFilterStr(req Req) {
	req.Respond(MatherialGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeMatherial(req Req) (Matherial, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var m Matherial
	err := decoder.Decode(&m)
	return m, err
}

func GetCash(req Req) {
	req.Respond(CashGet(req.IntParam, nil))
}

func GetCashAll(req Req) {
	req.Respond(CashGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateCash(req Req) {
	c, err := DecodeCash(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(CashCreate(c, nil))
}

func UpdateCash(req Req) {
	c, err := DecodeCash(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(CashUpdate(c, nil))
}

func DeleteCash(req Req) {
	req.Respond(CashDelete(req.IntParam, nil, false))
}

func GetCashByFilterInt(req Req) {
	req.Respond(CashGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetCashByFilterStr(req Req) {
	req.Respond(CashGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeCash(req Req) (Cash, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var c Cash
	err := decoder.Decode(&c)
	return c, err
}

func GetUserGroup(req Req) {
	req.Respond(UserGroupGet(req.IntParam, nil))
}

func GetUserGroupAll(req Req) {
	req.Respond(UserGroupGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateUserGroup(req Req) {
	u, err := DecodeUserGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(UserGroupCreate(u, nil))
}

func UpdateUserGroup(req Req) {
	u, err := DecodeUserGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(UserGroupUpdate(u, nil))
}

func DeleteUserGroup(req Req) {
	req.Respond(UserGroupDelete(req.IntParam, nil, false))
}

func GetUserGroupByFilterInt(req Req) {
	req.Respond(UserGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetUserGroupByFilterStr(req Req) {
	req.Respond(UserGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeUserGroup(req Req) (UserGroup, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var u UserGroup
	err := decoder.Decode(&u)
	return u, err
}

func GetUser(req Req) {
	req.Respond(UserGet(req.IntParam, nil))
}

func GetUserAll(req Req) {
	req.Respond(UserGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateUser(req Req) {
	u, err := DecodeUser(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(UserCreate(u, nil))
}

func UpdateUser(req Req) {
	u, err := DecodeUser(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(UserUpdate(u, nil))
}

func DeleteUser(req Req) {
	req.Respond(UserDelete(req.IntParam, nil, false))
}

func GetUserByFilterInt(req Req) {
	req.Respond(UserGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetUserByFilterStr(req Req) {
	req.Respond(UserGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeUser(req Req) (User, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var u User
	err := decoder.Decode(&u)
	return u, err
}

func GetEquipmentGroup(req Req) {
	req.Respond(EquipmentGroupGet(req.IntParam, nil))
}

func GetEquipmentGroupAll(req Req) {
	req.Respond(EquipmentGroupGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateEquipmentGroup(req Req) {
	e, err := DecodeEquipmentGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(EquipmentGroupCreate(e, nil))
}

func UpdateEquipmentGroup(req Req) {
	e, err := DecodeEquipmentGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(EquipmentGroupUpdate(e, nil))
}

func DeleteEquipmentGroup(req Req) {
	req.Respond(EquipmentGroupDelete(req.IntParam, nil, false))
}

func GetEquipmentGroupByFilterInt(req Req) {
	req.Respond(EquipmentGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetEquipmentGroupByFilterStr(req Req) {
	req.Respond(EquipmentGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeEquipmentGroup(req Req) (EquipmentGroup, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var e EquipmentGroup
	err := decoder.Decode(&e)
	return e, err
}

func GetEquipment(req Req) {
	req.Respond(EquipmentGet(req.IntParam, nil))
}

func GetEquipmentAll(req Req) {
	req.Respond(EquipmentGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateEquipment(req Req) {
	e, err := DecodeEquipment(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(EquipmentCreate(e, nil))
}

func UpdateEquipment(req Req) {
	e, err := DecodeEquipment(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(EquipmentUpdate(e, nil))
}

func DeleteEquipment(req Req) {
	req.Respond(EquipmentDelete(req.IntParam, nil, false))
}

func GetEquipmentByFilterInt(req Req) {
	req.Respond(EquipmentGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetEquipmentByFilterStr(req Req) {
	req.Respond(EquipmentGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeEquipment(req Req) (Equipment, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var e Equipment
	err := decoder.Decode(&e)
	return e, err
}

func GetOperationGroup(req Req) {
	req.Respond(OperationGroupGet(req.IntParam, nil))
}

func GetOperationGroupAll(req Req) {
	req.Respond(OperationGroupGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateOperationGroup(req Req) {
	o, err := DecodeOperationGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(OperationGroupCreate(o, nil))
}

func UpdateOperationGroup(req Req) {
	o, err := DecodeOperationGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(OperationGroupUpdate(o, nil))
}

func DeleteOperationGroup(req Req) {
	req.Respond(OperationGroupDelete(req.IntParam, nil, false))
}

func GetOperationGroupByFilterInt(req Req) {
	req.Respond(OperationGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetOperationGroupByFilterStr(req Req) {
	req.Respond(OperationGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeOperationGroup(req Req) (OperationGroup, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var o OperationGroup
	err := decoder.Decode(&o)
	return o, err
}

func GetOperation(req Req) {
	req.Respond(OperationGet(req.IntParam, nil))
}

func GetOperationAll(req Req) {
	req.Respond(OperationGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateOperation(req Req) {
	o, err := DecodeOperation(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(OperationCreate(o, nil))
}

func UpdateOperation(req Req) {
	o, err := DecodeOperation(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(OperationUpdate(o, nil))
}

func DeleteOperation(req Req) {
	req.Respond(OperationDelete(req.IntParam, nil, false))
}

func GetOperationByFilterInt(req Req) {
	req.Respond(OperationGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetOperationByFilterStr(req Req) {
	req.Respond(OperationGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeOperation(req Req) (Operation, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var o Operation
	err := decoder.Decode(&o)
	return o, err
}

func GetProductGroup(req Req) {
	req.Respond(ProductGroupGet(req.IntParam, nil))
}

func GetProductGroupAll(req Req) {
	req.Respond(ProductGroupGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateProductGroup(req Req) {
	p, err := DecodeProductGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProductGroupCreate(p, nil))
}

func UpdateProductGroup(req Req) {
	p, err := DecodeProductGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProductGroupUpdate(p, nil))
}

func DeleteProductGroup(req Req) {
	req.Respond(ProductGroupDelete(req.IntParam, nil, false))
}

func GetProductGroupByFilterInt(req Req) {
	req.Respond(ProductGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetProductGroupByFilterStr(req Req) {
	req.Respond(ProductGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeProductGroup(req Req) (ProductGroup, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var p ProductGroup
	err := decoder.Decode(&p)
	return p, err
}

func GetProduct(req Req) {
	req.Respond(ProductGet(req.IntParam, nil))
}

func GetProductAll(req Req) {
	req.Respond(ProductGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateProduct(req Req) {
	p, err := DecodeProduct(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProductCreate(p, nil))
}

func UpdateProduct(req Req) {
	p, err := DecodeProduct(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProductUpdate(p, nil))
}

func DeleteProduct(req Req) {
	req.Respond(ProductDelete(req.IntParam, nil, false))
}

func GetProductByFilterInt(req Req) {
	req.Respond(ProductGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetProductByFilterStr(req Req) {
	req.Respond(ProductGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeProduct(req Req) (Product, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var p Product
	err := decoder.Decode(&p)
	return p, err
}

func GetContragentGroup(req Req) {
	req.Respond(ContragentGroupGet(req.IntParam, nil))
}

func GetContragentGroupAll(req Req) {
	req.Respond(ContragentGroupGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateContragentGroup(req Req) {
	c, err := DecodeContragentGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ContragentGroupCreate(c, nil))
}

func UpdateContragentGroup(req Req) {
	c, err := DecodeContragentGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ContragentGroupUpdate(c, nil))
}

func DeleteContragentGroup(req Req) {
	req.Respond(ContragentGroupDelete(req.IntParam, nil, false))
}

func GetContragentGroupByFilterInt(req Req) {
	req.Respond(ContragentGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetContragentGroupByFilterStr(req Req) {
	req.Respond(ContragentGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeContragentGroup(req Req) (ContragentGroup, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var c ContragentGroup
	err := decoder.Decode(&c)
	return c, err
}

func GetContragent(req Req) {
	req.Respond(ContragentGet(req.IntParam, nil))
}

func GetContragentAll(req Req) {
	req.Respond(ContragentGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateContragent(req Req) {
	c, err := DecodeContragent(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ContragentCreate(c, nil))
}

func UpdateContragent(req Req) {
	c, err := DecodeContragent(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ContragentUpdate(c, nil))
}

func DeleteContragent(req Req) {
	req.Respond(ContragentDelete(req.IntParam, nil, false))
}

func GetContragentByFilterInt(req Req) {
	req.Respond(ContragentGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetContragentByFilterStr(req Req) {
	req.Respond(ContragentGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeContragent(req Req) (Contragent, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var c Contragent
	err := decoder.Decode(&c)
	return c, err
}

func GetContragentFindByContragentSearchContactSearch(req Req) {
	req.Respond(ContragentFindByContragentSearchContactSearch(req.StrParam))
}

func GetContact(req Req) {
	req.Respond(ContactGet(req.IntParam, nil))
}

func GetContactAll(req Req) {
	req.Respond(ContactGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateContact(req Req) {
	c, err := DecodeContact(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ContactCreate(c, nil))
}

func UpdateContact(req Req) {
	c, err := DecodeContact(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ContactUpdate(c, nil))
}

func DeleteContact(req Req) {
	req.Respond(ContactDelete(req.IntParam, nil, false))
}

func GetContactByFilterInt(req Req) {
	req.Respond(ContactGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetContactByFilterStr(req Req) {
	req.Respond(ContactGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeContact(req Req) (Contact, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var c Contact
	err := decoder.Decode(&c)
	return c, err
}

func GetOrderingStatus(req Req) {
	req.Respond(OrderingStatusGet(req.IntParam, nil))
}

func GetOrderingStatusAll(req Req) {
	req.Respond(OrderingStatusGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateOrderingStatus(req Req) {
	o, err := DecodeOrderingStatus(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(OrderingStatusCreate(o, nil))
}

func UpdateOrderingStatus(req Req) {
	o, err := DecodeOrderingStatus(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(OrderingStatusUpdate(o, nil))
}

func DeleteOrderingStatus(req Req) {
	req.Respond(OrderingStatusDelete(req.IntParam, nil, false))
}

func GetOrderingStatusByFilterInt(req Req) {
	req.Respond(OrderingStatusGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetOrderingStatusByFilterStr(req Req) {
	req.Respond(OrderingStatusGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeOrderingStatus(req Req) (OrderingStatus, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var o OrderingStatus
	err := decoder.Decode(&o)
	return o, err
}

func GetOrdering(req Req) {
	req.Respond(OrderingGet(req.IntParam, nil))
}

func GetOrderingAll(req Req) {
	req.Respond(OrderingGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateOrdering(req Req) {
	o, err := DecodeOrdering(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(OrderingCreate(o, nil))
}

func UpdateOrdering(req Req) {
	o, err := DecodeOrdering(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(OrderingUpdate(o, nil))
}

func UnRealizeOrdering(req Req) {
	req.Respond(OrderingDelete(req.IntParam, nil, true))
}

func DeleteOrdering(req Req) {
	req.Respond(OrderingDelete(req.IntParam, nil, false))
}

func GetOrderingByFilterInt(req Req) {
	req.Respond(OrderingGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetOrderingByFilterStr(req Req) {
	req.Respond(OrderingGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeOrdering(req Req) (Ordering, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var o Ordering
	err := decoder.Decode(&o)
	return o, err
}

func GetOrderingBetweenCreatedAt(req Req) {
	req.Respond(OrderingGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetOrderingBetweenDeadlineAt(req Req) {
	req.Respond(OrderingGetBetweenDeadlineAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetOrderingCashSumSumBefore(req Req) {
	req.Respond(OrderingCashSumGetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}

func GetOrderingSumByFilter(req Req) {
	req.Respond(OrderingGetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}

func GetOwner(req Req) {
	req.Respond(OwnerGet(req.IntParam, nil))
}

func GetOwnerAll(req Req) {
	req.Respond(OwnerGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateOwner(req Req) {
	o, err := DecodeOwner(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(OwnerCreate(o, nil))
}

func UpdateOwner(req Req) {
	o, err := DecodeOwner(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(OwnerUpdate(o, nil))
}

func DeleteOwner(req Req) {
	req.Respond(OwnerDelete(req.IntParam, nil, false))
}

func GetOwnerByFilterInt(req Req) {
	req.Respond(OwnerGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetOwnerByFilterStr(req Req) {
	req.Respond(OwnerGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeOwner(req Req) (Owner, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var o Owner
	err := decoder.Decode(&o)
	return o, err
}

func GetInvoice(req Req) {
	req.Respond(InvoiceGet(req.IntParam, nil))
}

func GetInvoiceAll(req Req) {
	req.Respond(InvoiceGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateInvoice(req Req) {
	i, err := DecodeInvoice(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(InvoiceCreate(i, nil))
}

func UpdateInvoice(req Req) {
	i, err := DecodeInvoice(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(InvoiceUpdate(i, nil))
}

func UnRealizeInvoice(req Req) {
	req.Respond(InvoiceDelete(req.IntParam, nil, true))
}

func DeleteInvoice(req Req) {
	req.Respond(InvoiceDelete(req.IntParam, nil, false))
}

func GetInvoiceByFilterInt(req Req) {
	req.Respond(InvoiceGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetInvoiceByFilterStr(req Req) {
	req.Respond(InvoiceGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeInvoice(req Req) (Invoice, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var i Invoice
	err := decoder.Decode(&i)
	return i, err
}

func RealizedInvoice(req Req) {
	req.Respond(InvoiceRealized(req.IntParam, nil))
}

func GetInvoiceBetweenCreatedAt(req Req) {
	req.Respond(InvoiceGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetInvoiceCashSumSumBefore(req Req) {
	req.Respond(InvoiceCashSumGetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}

func GetInvoiceSumByFilter(req Req) {
	req.Respond(InvoiceGetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}

func GetItemToInvoice(req Req) {
	req.Respond(ItemToInvoiceGet(req.IntParam, nil))
}

func GetItemToInvoiceAll(req Req) {
	req.Respond(ItemToInvoiceGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateItemToInvoice(req Req) {
	i, err := DecodeItemToInvoice(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ItemToInvoiceCreate(i, nil))
}

func UpdateItemToInvoice(req Req) {
	i, err := DecodeItemToInvoice(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ItemToInvoiceUpdate(i, nil))
}

func DeleteItemToInvoice(req Req) {
	req.Respond(ItemToInvoiceDelete(req.IntParam, nil, false))
}

func GetItemToInvoiceByFilterInt(req Req) {
	req.Respond(ItemToInvoiceGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetItemToInvoiceByFilterStr(req Req) {
	req.Respond(ItemToInvoiceGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeItemToInvoice(req Req) (ItemToInvoice, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var i ItemToInvoice
	err := decoder.Decode(&i)
	return i, err
}

func GetItemToInvoiceCostSumBefore(req Req) {
	req.Respond(ItemToInvoiceCostGetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}

func GetItemToInvoiceSumByFilter(req Req) {
	req.Respond(ItemToInvoiceGetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}

func GetProductToOrderingStatus(req Req) {
	req.Respond(ProductToOrderingStatusGet(req.IntParam, nil))
}

func GetProductToOrderingStatusAll(req Req) {
	req.Respond(ProductToOrderingStatusGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateProductToOrderingStatus(req Req) {
	p, err := DecodeProductToOrderingStatus(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProductToOrderingStatusCreate(p, nil))
}

func UpdateProductToOrderingStatus(req Req) {
	p, err := DecodeProductToOrderingStatus(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProductToOrderingStatusUpdate(p, nil))
}

func DeleteProductToOrderingStatus(req Req) {
	req.Respond(ProductToOrderingStatusDelete(req.IntParam, nil, false))
}

func GetProductToOrderingStatusByFilterInt(req Req) {
	req.Respond(ProductToOrderingStatusGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetProductToOrderingStatusByFilterStr(req Req) {
	req.Respond(ProductToOrderingStatusGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeProductToOrderingStatus(req Req) (ProductToOrderingStatus, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var p ProductToOrderingStatus
	err := decoder.Decode(&p)
	return p, err
}

func GetProductToOrdering(req Req) {
	req.Respond(ProductToOrderingGet(req.IntParam, nil))
}

func GetProductToOrderingAll(req Req) {
	req.Respond(ProductToOrderingGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateProductToOrdering(req Req) {
	p, err := DecodeProductToOrdering(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProductToOrderingCreate(p, nil))
}

func UpdateProductToOrdering(req Req) {
	p, err := DecodeProductToOrdering(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProductToOrderingUpdate(p, nil))
}

func DeleteProductToOrdering(req Req) {
	req.Respond(ProductToOrderingDelete(req.IntParam, nil, false))
}

func GetProductToOrderingByFilterInt(req Req) {
	req.Respond(ProductToOrderingGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetProductToOrderingByFilterStr(req Req) {
	req.Respond(ProductToOrderingGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeProductToOrdering(req Req) (ProductToOrdering, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var p ProductToOrdering
	err := decoder.Decode(&p)
	return p, err
}

func GetProductToOrderingCostSumBefore(req Req) {
	req.Respond(ProductToOrderingCostGetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}

func GetProductToOrderingSumByFilter(req Req) {
	req.Respond(ProductToOrderingGetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}

func GetMatherialToOrdering(req Req) {
	req.Respond(MatherialToOrderingGet(req.IntParam, nil))
}

func GetMatherialToOrderingAll(req Req) {
	req.Respond(MatherialToOrderingGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateMatherialToOrdering(req Req) {
	m, err := DecodeMatherialToOrdering(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialToOrderingCreate(m, nil))
}

func UpdateMatherialToOrdering(req Req) {
	m, err := DecodeMatherialToOrdering(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialToOrderingUpdate(m, nil))
}

func DeleteMatherialToOrdering(req Req) {
	req.Respond(MatherialToOrderingDelete(req.IntParam, nil, false))
}

func GetMatherialToOrderingByFilterInt(req Req) {
	req.Respond(MatherialToOrderingGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetMatherialToOrderingByFilterStr(req Req) {
	req.Respond(MatherialToOrderingGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeMatherialToOrdering(req Req) (MatherialToOrdering, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var m MatherialToOrdering
	err := decoder.Decode(&m)
	return m, err
}

func GetMatherialToOrderingCostSumBefore(req Req) {
	req.Respond(MatherialToOrderingCostGetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}

func GetMatherialToOrderingSumByFilter(req Req) {
	req.Respond(MatherialToOrderingGetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}

func GetMatherialToProduct(req Req) {
	req.Respond(MatherialToProductGet(req.IntParam, nil))
}

func GetMatherialToProductAll(req Req) {
	req.Respond(MatherialToProductGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateMatherialToProduct(req Req) {
	m, err := DecodeMatherialToProduct(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialToProductCreate(m, nil))
}

func UpdateMatherialToProduct(req Req) {
	m, err := DecodeMatherialToProduct(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialToProductUpdate(m, nil))
}

func DeleteMatherialToProduct(req Req) {
	req.Respond(MatherialToProductDelete(req.IntParam, nil, false))
}

func GetMatherialToProductByFilterInt(req Req) {
	req.Respond(MatherialToProductGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetMatherialToProductByFilterStr(req Req) {
	req.Respond(MatherialToProductGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeMatherialToProduct(req Req) (MatherialToProduct, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var m MatherialToProduct
	err := decoder.Decode(&m)
	return m, err
}

func GetMatherialToProductCostSumBefore(req Req) {
	req.Respond(MatherialToProductCostGetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}

func GetMatherialToProductSumByFilter(req Req) {
	req.Respond(MatherialToProductGetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}

func GetOperationToOrdering(req Req) {
	req.Respond(OperationToOrderingGet(req.IntParam, nil))
}

func GetOperationToOrderingAll(req Req) {
	req.Respond(OperationToOrderingGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateOperationToOrdering(req Req) {
	o, err := DecodeOperationToOrdering(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(OperationToOrderingCreate(o, nil))
}

func UpdateOperationToOrdering(req Req) {
	o, err := DecodeOperationToOrdering(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(OperationToOrderingUpdate(o, nil))
}

func DeleteOperationToOrdering(req Req) {
	req.Respond(OperationToOrderingDelete(req.IntParam, nil, false))
}

func GetOperationToOrderingByFilterInt(req Req) {
	req.Respond(OperationToOrderingGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetOperationToOrderingByFilterStr(req Req) {
	req.Respond(OperationToOrderingGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeOperationToOrdering(req Req) (OperationToOrdering, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var o OperationToOrdering
	err := decoder.Decode(&o)
	return o, err
}

func GetOperationToOrderingCostSumBefore(req Req) {
	req.Respond(OperationToOrderingCostGetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}

func GetOperationToOrderingSumByFilter(req Req) {
	req.Respond(OperationToOrderingGetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}

func GetOperationToProduct(req Req) {
	req.Respond(OperationToProductGet(req.IntParam, nil))
}

func GetOperationToProductAll(req Req) {
	req.Respond(OperationToProductGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateOperationToProduct(req Req) {
	o, err := DecodeOperationToProduct(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(OperationToProductCreate(o, nil))
}

func UpdateOperationToProduct(req Req) {
	o, err := DecodeOperationToProduct(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(OperationToProductUpdate(o, nil))
}

func DeleteOperationToProduct(req Req) {
	req.Respond(OperationToProductDelete(req.IntParam, nil, false))
}

func GetOperationToProductByFilterInt(req Req) {
	req.Respond(OperationToProductGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetOperationToProductByFilterStr(req Req) {
	req.Respond(OperationToProductGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeOperationToProduct(req Req) (OperationToProduct, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var o OperationToProduct
	err := decoder.Decode(&o)
	return o, err
}

func GetOperationToProductCostSumBefore(req Req) {
	req.Respond(OperationToProductCostGetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}

func GetOperationToProductSumByFilter(req Req) {
	req.Respond(OperationToProductGetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}

func GetProductToProduct(req Req) {
	req.Respond(ProductToProductGet(req.IntParam, nil))
}

func GetProductToProductAll(req Req) {
	req.Respond(ProductToProductGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateProductToProduct(req Req) {
	p, err := DecodeProductToProduct(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProductToProductCreate(p, nil))
}

func UpdateProductToProduct(req Req) {
	p, err := DecodeProductToProduct(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProductToProductUpdate(p, nil))
}

func DeleteProductToProduct(req Req) {
	req.Respond(ProductToProductDelete(req.IntParam, nil, false))
}

func GetProductToProductByFilterInt(req Req) {
	req.Respond(ProductToProductGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetProductToProductByFilterStr(req Req) {
	req.Respond(ProductToProductGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeProductToProduct(req Req) (ProductToProduct, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var p ProductToProduct
	err := decoder.Decode(&p)
	return p, err
}

func GetProductToProductCostSumBefore(req Req) {
	req.Respond(ProductToProductCostGetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}

func GetProductToProductSumByFilter(req Req) {
	req.Respond(ProductToProductGetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}

func GetCboxCheck(req Req) {
	req.Respond(CboxCheckGet(req.IntParam, nil))
}

func GetCboxCheckAll(req Req) {
	req.Respond(CboxCheckGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateCboxCheck(req Req) {
	c, err := DecodeCboxCheck(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(CboxCheckCreate(c, nil))
}

func UpdateCboxCheck(req Req) {
	c, err := DecodeCboxCheck(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(CboxCheckUpdate(c, nil))
}

func DeleteCboxCheck(req Req) {
	req.Respond(CboxCheckDelete(req.IntParam, nil, false))
}

func GetCboxCheckByFilterInt(req Req) {
	req.Respond(CboxCheckGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetCboxCheckByFilterStr(req Req) {
	req.Respond(CboxCheckGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeCboxCheck(req Req) (CboxCheck, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var c CboxCheck
	err := decoder.Decode(&c)
	return c, err
}

func GetCboxCheckBetweenCreatedAt(req Req) {
	req.Respond(CboxCheckGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetItemToCboxCheck(req Req) {
	req.Respond(ItemToCboxCheckGet(req.IntParam, nil))
}

func GetItemToCboxCheckAll(req Req) {
	req.Respond(ItemToCboxCheckGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateItemToCboxCheck(req Req) {
	i, err := DecodeItemToCboxCheck(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ItemToCboxCheckCreate(i, nil))
}

func UpdateItemToCboxCheck(req Req) {
	i, err := DecodeItemToCboxCheck(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ItemToCboxCheckUpdate(i, nil))
}

func DeleteItemToCboxCheck(req Req) {
	req.Respond(ItemToCboxCheckDelete(req.IntParam, nil, false))
}

func GetItemToCboxCheckByFilterInt(req Req) {
	req.Respond(ItemToCboxCheckGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetItemToCboxCheckByFilterStr(req Req) {
	req.Respond(ItemToCboxCheckGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeItemToCboxCheck(req Req) (ItemToCboxCheck, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var i ItemToCboxCheck
	err := decoder.Decode(&i)
	return i, err
}

func GetItemToCboxCheckCostSumBefore(req Req) {
	req.Respond(ItemToCboxCheckCostGetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}

func GetItemToCboxCheckSumByFilter(req Req) {
	req.Respond(ItemToCboxCheckGetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}

func GetCashIn(req Req) {
	req.Respond(CashInGet(req.IntParam, nil))
}

func GetCashInAll(req Req) {
	req.Respond(CashInGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateCashIn(req Req) {
	c, err := DecodeCashIn(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(CashInCreate(c, nil))
}

func UpdateCashIn(req Req) {
	c, err := DecodeCashIn(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(CashInUpdate(c, nil))
}

func UnRealizeCashIn(req Req) {
	req.Respond(CashInDelete(req.IntParam, nil, true))
}

func DeleteCashIn(req Req) {
	req.Respond(CashInDelete(req.IntParam, nil, false))
}

func GetCashInByFilterInt(req Req) {
	req.Respond(CashInGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetCashInByFilterStr(req Req) {
	req.Respond(CashInGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeCashIn(req Req) (CashIn, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var c CashIn
	err := decoder.Decode(&c)
	return c, err
}

func RealizedCashIn(req Req) {
	req.Respond(CashInRealized(req.IntParam, nil))
}

func GetCashInBetweenCreatedAt(req Req) {
	req.Respond(CashInGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetCashInCashSumSumBefore(req Req) {
	req.Respond(CashInCashSumGetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}

func GetCashInSumByFilter(req Req) {
	req.Respond(CashInGetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}

func GetCashOut(req Req) {
	req.Respond(CashOutGet(req.IntParam, nil))
}

func GetCashOutAll(req Req) {
	req.Respond(CashOutGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateCashOut(req Req) {
	c, err := DecodeCashOut(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(CashOutCreate(c, nil))
}

func UpdateCashOut(req Req) {
	c, err := DecodeCashOut(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(CashOutUpdate(c, nil))
}

func UnRealizeCashOut(req Req) {
	req.Respond(CashOutDelete(req.IntParam, nil, true))
}

func DeleteCashOut(req Req) {
	req.Respond(CashOutDelete(req.IntParam, nil, false))
}

func GetCashOutByFilterInt(req Req) {
	req.Respond(CashOutGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetCashOutByFilterStr(req Req) {
	req.Respond(CashOutGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeCashOut(req Req) (CashOut, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var c CashOut
	err := decoder.Decode(&c)
	return c, err
}

func RealizedCashOut(req Req) {
	req.Respond(CashOutRealized(req.IntParam, nil))
}

func GetCashOutBetweenCreatedAt(req Req) {
	req.Respond(CashOutGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetCashOutCashSumSumBefore(req Req) {
	req.Respond(CashOutCashSumGetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}

func GetCashOutSumByFilter(req Req) {
	req.Respond(CashOutGetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}

func GetWhs(req Req) {
	req.Respond(WhsGet(req.IntParam, nil))
}

func GetWhsAll(req Req) {
	req.Respond(WhsGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateWhs(req Req) {
	w, err := DecodeWhs(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(WhsCreate(w, nil))
}

func UpdateWhs(req Req) {
	w, err := DecodeWhs(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(WhsUpdate(w, nil))
}

func DeleteWhs(req Req) {
	req.Respond(WhsDelete(req.IntParam, nil, false))
}

func GetWhsByFilterInt(req Req) {
	req.Respond(WhsGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetWhsByFilterStr(req Req) {
	req.Respond(WhsGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeWhs(req Req) (Whs, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var w Whs
	err := decoder.Decode(&w)
	return w, err
}

func GetWhsIn(req Req) {
	req.Respond(WhsInGet(req.IntParam, nil))
}

func GetWhsInAll(req Req) {
	req.Respond(WhsInGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateWhsIn(req Req) {
	w, err := DecodeWhsIn(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(WhsInCreate(w, nil))
}

func UpdateWhsIn(req Req) {
	w, err := DecodeWhsIn(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(WhsInUpdate(w, nil))
}

func UnRealizeWhsIn(req Req) {
	req.Respond(WhsInDelete(req.IntParam, nil, true))
}

func DeleteWhsIn(req Req) {
	req.Respond(WhsInDelete(req.IntParam, nil, false))
}

func GetWhsInByFilterInt(req Req) {
	req.Respond(WhsInGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetWhsInByFilterStr(req Req) {
	req.Respond(WhsInGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeWhsIn(req Req) (WhsIn, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var w WhsIn
	err := decoder.Decode(&w)
	return w, err
}

func RealizedWhsIn(req Req) {
	req.Respond(WhsInRealized(req.IntParam, nil))
}

func GetWhsInBetweenCreatedAt(req Req) {
	req.Respond(WhsInGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWhsInBetweenContragentCreatedAt(req Req) {
	req.Respond(WhsInGetBetweenContragentCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWhsInWhsSumSumBefore(req Req) {
	req.Respond(WhsInWhsSumGetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}

func GetWhsInSumByFilter(req Req) {
	req.Respond(WhsInGetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}

func GetWhsOut(req Req) {
	req.Respond(WhsOutGet(req.IntParam, nil))
}

func GetWhsOutAll(req Req) {
	req.Respond(WhsOutGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateWhsOut(req Req) {
	w, err := DecodeWhsOut(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(WhsOutCreate(w, nil))
}

func UpdateWhsOut(req Req) {
	w, err := DecodeWhsOut(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(WhsOutUpdate(w, nil))
}

func UnRealizeWhsOut(req Req) {
	req.Respond(WhsOutDelete(req.IntParam, nil, true))
}

func DeleteWhsOut(req Req) {
	req.Respond(WhsOutDelete(req.IntParam, nil, false))
}

func GetWhsOutByFilterInt(req Req) {
	req.Respond(WhsOutGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetWhsOutByFilterStr(req Req) {
	req.Respond(WhsOutGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeWhsOut(req Req) (WhsOut, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var w WhsOut
	err := decoder.Decode(&w)
	return w, err
}

func RealizedWhsOut(req Req) {
	req.Respond(WhsOutRealized(req.IntParam, nil))
}

func GetWhsOutBetweenCreatedAt(req Req) {
	req.Respond(WhsOutGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWhsOutWhsSumSumBefore(req Req) {
	req.Respond(WhsOutWhsSumGetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}

func GetWhsOutSumByFilter(req Req) {
	req.Respond(WhsOutGetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}

func GetMatherialToWhsIn(req Req) {
	req.Respond(MatherialToWhsInGet(req.IntParam, nil))
}

func GetMatherialToWhsInAll(req Req) {
	req.Respond(MatherialToWhsInGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateMatherialToWhsIn(req Req) {
	m, err := DecodeMatherialToWhsIn(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialToWhsInCreate(m, nil))
}

func UpdateMatherialToWhsIn(req Req) {
	m, err := DecodeMatherialToWhsIn(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialToWhsInUpdate(m, nil))
}

func DeleteMatherialToWhsIn(req Req) {
	req.Respond(MatherialToWhsInDelete(req.IntParam, nil, false))
}

func GetMatherialToWhsInByFilterInt(req Req) {
	req.Respond(MatherialToWhsInGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetMatherialToWhsInByFilterStr(req Req) {
	req.Respond(MatherialToWhsInGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeMatherialToWhsIn(req Req) (MatherialToWhsIn, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var m MatherialToWhsIn
	err := decoder.Decode(&m)
	return m, err
}

func GetMatherialToWhsOut(req Req) {
	req.Respond(MatherialToWhsOutGet(req.IntParam, nil))
}

func GetMatherialToWhsOutAll(req Req) {
	req.Respond(MatherialToWhsOutGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateMatherialToWhsOut(req Req) {
	m, err := DecodeMatherialToWhsOut(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialToWhsOutCreate(m, nil))
}

func UpdateMatherialToWhsOut(req Req) {
	m, err := DecodeMatherialToWhsOut(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialToWhsOutUpdate(m, nil))
}

func DeleteMatherialToWhsOut(req Req) {
	req.Respond(MatherialToWhsOutDelete(req.IntParam, nil, false))
}

func GetMatherialToWhsOutByFilterInt(req Req) {
	req.Respond(MatherialToWhsOutGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetMatherialToWhsOutByFilterStr(req Req) {
	req.Respond(MatherialToWhsOutGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeMatherialToWhsOut(req Req) (MatherialToWhsOut, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var m MatherialToWhsOut
	err := decoder.Decode(&m)
	return m, err
}

func GetMatherialPart(req Req) {
	req.Respond(MatherialPartGet(req.IntParam, nil))
}

func GetMatherialPartAll(req Req) {
	req.Respond(MatherialPartGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateMatherialPart(req Req) {
	m, err := DecodeMatherialPart(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialPartCreate(m, nil))
}

func UpdateMatherialPart(req Req) {
	m, err := DecodeMatherialPart(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialPartUpdate(m, nil))
}

func DeleteMatherialPart(req Req) {
	req.Respond(MatherialPartDelete(req.IntParam, nil, false))
}

func GetMatherialPartByFilterInt(req Req) {
	req.Respond(MatherialPartGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetMatherialPartByFilterStr(req Req) {
	req.Respond(MatherialPartGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeMatherialPart(req Req) (MatherialPart, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var m MatherialPart
	err := decoder.Decode(&m)
	return m, err
}

func GetMatherialPartBetweenCreatedAt(req Req) {
	req.Respond(MatherialPartGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetMatherialPartSlice(req Req) {
	req.Respond(MatherialPartSliceGet(req.IntParam, nil))
}

func GetMatherialPartSliceAll(req Req) {
	req.Respond(MatherialPartSliceGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateMatherialPartSlice(req Req) {
	m, err := DecodeMatherialPartSlice(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialPartSliceCreate(m, nil))
}

func UpdateMatherialPartSlice(req Req) {
	m, err := DecodeMatherialPartSlice(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(MatherialPartSliceUpdate(m, nil))
}

func DeleteMatherialPartSlice(req Req) {
	req.Respond(MatherialPartSliceDelete(req.IntParam, nil, false))
}

func GetMatherialPartSliceByFilterInt(req Req) {
	req.Respond(MatherialPartSliceGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetMatherialPartSliceByFilterStr(req Req) {
	req.Respond(MatherialPartSliceGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeMatherialPartSlice(req Req) (MatherialPartSlice, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var m MatherialPartSlice
	err := decoder.Decode(&m)
	return m, err
}

func GetMatherialPartSliceBetweenCreatedAt(req Req) {
	req.Respond(MatherialPartSliceGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetProjectGroup(req Req) {
	req.Respond(ProjectGroupGet(req.IntParam, nil))
}

func GetProjectGroupAll(req Req) {
	req.Respond(ProjectGroupGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateProjectGroup(req Req) {
	p, err := DecodeProjectGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProjectGroupCreate(p, nil))
}

func UpdateProjectGroup(req Req) {
	p, err := DecodeProjectGroup(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProjectGroupUpdate(p, nil))
}

func DeleteProjectGroup(req Req) {
	req.Respond(ProjectGroupDelete(req.IntParam, nil, false))
}

func GetProjectGroupByFilterInt(req Req) {
	req.Respond(ProjectGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetProjectGroupByFilterStr(req Req) {
	req.Respond(ProjectGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeProjectGroup(req Req) (ProjectGroup, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var p ProjectGroup
	err := decoder.Decode(&p)
	return p, err
}

func GetProjectStatus(req Req) {
	req.Respond(ProjectStatusGet(req.IntParam, nil))
}

func GetProjectStatusAll(req Req) {
	req.Respond(ProjectStatusGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateProjectStatus(req Req) {
	p, err := DecodeProjectStatus(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProjectStatusCreate(p, nil))
}

func UpdateProjectStatus(req Req) {
	p, err := DecodeProjectStatus(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProjectStatusUpdate(p, nil))
}

func DeleteProjectStatus(req Req) {
	req.Respond(ProjectStatusDelete(req.IntParam, nil, false))
}

func GetProjectStatusByFilterInt(req Req) {
	req.Respond(ProjectStatusGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetProjectStatusByFilterStr(req Req) {
	req.Respond(ProjectStatusGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeProjectStatus(req Req) (ProjectStatus, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var p ProjectStatus
	err := decoder.Decode(&p)
	return p, err
}

func GetProjectType(req Req) {
	req.Respond(ProjectTypeGet(req.IntParam, nil))
}

func GetProjectTypeAll(req Req) {
	req.Respond(ProjectTypeGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateProjectType(req Req) {
	p, err := DecodeProjectType(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProjectTypeCreate(p, nil))
}

func UpdateProjectType(req Req) {
	p, err := DecodeProjectType(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProjectTypeUpdate(p, nil))
}

func DeleteProjectType(req Req) {
	req.Respond(ProjectTypeDelete(req.IntParam, nil, false))
}

func GetProjectTypeByFilterInt(req Req) {
	req.Respond(ProjectTypeGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetProjectTypeByFilterStr(req Req) {
	req.Respond(ProjectTypeGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeProjectType(req Req) (ProjectType, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var p ProjectType
	err := decoder.Decode(&p)
	return p, err
}

func GetProject(req Req) {
	req.Respond(ProjectGet(req.IntParam, nil))
}

func GetProjectAll(req Req) {
	req.Respond(ProjectGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateProject(req Req) {
	p, err := DecodeProject(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProjectCreate(p, nil))
}

func UpdateProject(req Req) {
	p, err := DecodeProject(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(ProjectUpdate(p, nil))
}

func DeleteProject(req Req) {
	req.Respond(ProjectDelete(req.IntParam, nil, false))
}

func GetProjectByFilterInt(req Req) {
	req.Respond(ProjectGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetProjectByFilterStr(req Req) {
	req.Respond(ProjectGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeProject(req Req) (Project, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var p Project
	err := decoder.Decode(&p)
	return p, err
}

func GetProjectFindByProjectInfoContragentNoSearchContactNoSearch(req Req) {
	req.Respond(ProjectFindByProjectInfoContragentNoSearchContactNoSearch(req.StrParam))
}

func GetProjectBetweenCreatedAt(req Req) {
	req.Respond(ProjectGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetCounter(req Req) {
	req.Respond(CounterGet(req.IntParam, nil))
}

func GetCounterAll(req Req) {
	req.Respond(CounterGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateCounter(req Req) {
	c, err := DecodeCounter(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(CounterCreate(c, nil))
}

func UpdateCounter(req Req) {
	c, err := DecodeCounter(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(CounterUpdate(c, nil))
}

func DeleteCounter(req Req) {
	req.Respond(CounterDelete(req.IntParam, nil, false))
}

func GetCounterByFilterInt(req Req) {
	req.Respond(CounterGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetCounterByFilterStr(req Req) {
	req.Respond(CounterGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeCounter(req Req) (Counter, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var c Counter
	err := decoder.Decode(&c)
	return c, err
}

func GetRecordToCounter(req Req) {
	req.Respond(RecordToCounterGet(req.IntParam, nil))
}

func GetRecordToCounterAll(req Req) {
	req.Respond(RecordToCounterGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateRecordToCounter(req Req) {
	r, err := DecodeRecordToCounter(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(RecordToCounterCreate(r, nil))
}

func UpdateRecordToCounter(req Req) {
	r, err := DecodeRecordToCounter(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(RecordToCounterUpdate(r, nil))
}

func GetRecordToCounterByFilterInt(req Req) {
	req.Respond(RecordToCounterGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetRecordToCounterByFilterStr(req Req) {
	req.Respond(RecordToCounterGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeRecordToCounter(req Req) (RecordToCounter, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var r RecordToCounter
	err := decoder.Decode(&r)
	return r, err
}

func GetRecordToCounterBetweenCreatedAt(req Req) {
	req.Respond(RecordToCounterGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetRecordToCounterNumberSumBefore(req Req) {
	req.Respond(RecordToCounterNumberGetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}

func GetRecordToCounterSumByFilter(req Req) {
	req.Respond(RecordToCounterGetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}

func GetWmcNumber(req Req) {
	req.Respond(WmcNumberGet(req.IntParam, nil))
}

func GetWmcNumberAll(req Req) {
	req.Respond(WmcNumberGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateWmcNumber(req Req) {
	w, err := DecodeWmcNumber(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(WmcNumberCreate(w, nil))
}

func UpdateWmcNumber(req Req) {
	w, err := DecodeWmcNumber(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(WmcNumberUpdate(w, nil))
}

func DeleteWmcNumber(req Req) {
	req.Respond(WmcNumberDelete(req.IntParam, nil, false))
}

func GetWmcNumberByFilterInt(req Req) {
	req.Respond(WmcNumberGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetWmcNumberByFilterStr(req Req) {
	req.Respond(WmcNumberGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeWmcNumber(req Req) (WmcNumber, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var w WmcNumber
	err := decoder.Decode(&w)
	return w, err
}

func GetNumbersToProduct(req Req) {
	req.Respond(NumbersToProductGet(req.IntParam, nil))
}

func GetNumbersToProductAll(req Req) {
	req.Respond(NumbersToProductGetAll(req.WithDeleted, req.DeletedOnly, nil))
}

func CreateNumbersToProduct(req Req) {
	n, err := DecodeNumbersToProduct(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(NumbersToProductCreate(n, nil))
}

func UpdateNumbersToProduct(req Req) {
	n, err := DecodeNumbersToProduct(req)
	if err != nil {
		req.Respond(nil, err)
		return
	}
	req.Respond(NumbersToProductUpdate(n, nil))
}

func DeleteNumbersToProduct(req Req) {
	req.Respond(NumbersToProductDelete(req.IntParam, nil, false))
}

func GetNumbersToProductByFilterInt(req Req) {
	req.Respond(NumbersToProductGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
}

func GetNumbersToProductByFilterStr(req Req) {
	req.Respond(NumbersToProductGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
}

func DecodeNumbersToProduct(req Req) (NumbersToProduct, error) {
	decoder := json.NewDecoder(req.R.Body)
	defer req.R.Body.Close()
	var n NumbersToProduct
	err := decoder.Decode(&n)
	return n, err
}

func GetWDocument(req Req) {
	req.Respond(WDocumentGet(req.IntParam))
}

func GetWDocumentAll(req Req) {
	req.Respond(WDocumentGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWDocumentByFilterInt(req Req) {
	req.Respond(WDocumentGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWDocumentByFilterStr(req Req) {
	req.Respond(WDocumentGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWMeasure(req Req) {
	req.Respond(WMeasureGet(req.IntParam))
}

func GetWMeasureAll(req Req) {
	req.Respond(WMeasureGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWMeasureByFilterInt(req Req) {
	req.Respond(WMeasureGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWMeasureByFilterStr(req Req) {
	req.Respond(WMeasureGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWCountType(req Req) {
	req.Respond(WCountTypeGet(req.IntParam))
}

func GetWCountTypeAll(req Req) {
	req.Respond(WCountTypeGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWCountTypeByFilterInt(req Req) {
	req.Respond(WCountTypeGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWCountTypeByFilterStr(req Req) {
	req.Respond(WCountTypeGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWColorGroup(req Req) {
	req.Respond(WColorGroupGet(req.IntParam))
}

func GetWColorGroupAll(req Req) {
	req.Respond(WColorGroupGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWColorGroupByFilterInt(req Req) {
	req.Respond(WColorGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWColorGroupByFilterStr(req Req) {
	req.Respond(WColorGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWColor(req Req) {
	req.Respond(WColorGet(req.IntParam))
}

func GetWColorAll(req Req) {
	req.Respond(WColorGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWColorByFilterInt(req Req) {
	req.Respond(WColorGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWColorByFilterStr(req Req) {
	req.Respond(WColorGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialGroup(req Req) {
	req.Respond(WMatherialGroupGet(req.IntParam))
}

func GetWMatherialGroupAll(req Req) {
	req.Respond(WMatherialGroupGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialGroupByFilterInt(req Req) {
	req.Respond(WMatherialGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialGroupByFilterStr(req Req) {
	req.Respond(WMatherialGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherial(req Req) {
	req.Respond(WMatherialGet(req.IntParam))
}

func GetWMatherialAll(req Req) {
	req.Respond(WMatherialGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialByFilterInt(req Req) {
	req.Respond(WMatherialGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialByFilterStr(req Req) {
	req.Respond(WMatherialGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWCash(req Req) {
	req.Respond(WCashGet(req.IntParam))
}

func GetWCashAll(req Req) {
	req.Respond(WCashGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWCashByFilterInt(req Req) {
	req.Respond(WCashGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWCashByFilterStr(req Req) {
	req.Respond(WCashGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWUserGroup(req Req) {
	req.Respond(WUserGroupGet(req.IntParam))
}

func GetWUserGroupAll(req Req) {
	req.Respond(WUserGroupGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWUserGroupByFilterInt(req Req) {
	req.Respond(WUserGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWUserGroupByFilterStr(req Req) {
	req.Respond(WUserGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWUser(req Req) {
	req.Respond(WUserGet(req.IntParam))
}

func GetWUserAll(req Req) {
	req.Respond(WUserGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWUserByFilterInt(req Req) {
	req.Respond(WUserGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWUserByFilterStr(req Req) {
	req.Respond(WUserGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWEquipmentGroup(req Req) {
	req.Respond(WEquipmentGroupGet(req.IntParam))
}

func GetWEquipmentGroupAll(req Req) {
	req.Respond(WEquipmentGroupGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWEquipmentGroupByFilterInt(req Req) {
	req.Respond(WEquipmentGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWEquipmentGroupByFilterStr(req Req) {
	req.Respond(WEquipmentGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWEquipment(req Req) {
	req.Respond(WEquipmentGet(req.IntParam))
}

func GetWEquipmentAll(req Req) {
	req.Respond(WEquipmentGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWEquipmentByFilterInt(req Req) {
	req.Respond(WEquipmentGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWEquipmentByFilterStr(req Req) {
	req.Respond(WEquipmentGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWOperationGroup(req Req) {
	req.Respond(WOperationGroupGet(req.IntParam))
}

func GetWOperationGroupAll(req Req) {
	req.Respond(WOperationGroupGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWOperationGroupByFilterInt(req Req) {
	req.Respond(WOperationGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWOperationGroupByFilterStr(req Req) {
	req.Respond(WOperationGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWOperation(req Req) {
	req.Respond(WOperationGet(req.IntParam))
}

func GetWOperationAll(req Req) {
	req.Respond(WOperationGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWOperationByFilterInt(req Req) {
	req.Respond(WOperationGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWOperationByFilterStr(req Req) {
	req.Respond(WOperationGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWProductGroup(req Req) {
	req.Respond(WProductGroupGet(req.IntParam))
}

func GetWProductGroupAll(req Req) {
	req.Respond(WProductGroupGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWProductGroupByFilterInt(req Req) {
	req.Respond(WProductGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWProductGroupByFilterStr(req Req) {
	req.Respond(WProductGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWProduct(req Req) {
	req.Respond(WProductGet(req.IntParam))
}

func GetWProductAll(req Req) {
	req.Respond(WProductGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWProductByFilterInt(req Req) {
	req.Respond(WProductGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWProductByFilterStr(req Req) {
	req.Respond(WProductGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWContragentGroup(req Req) {
	req.Respond(WContragentGroupGet(req.IntParam))
}

func GetWContragentGroupAll(req Req) {
	req.Respond(WContragentGroupGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWContragentGroupByFilterInt(req Req) {
	req.Respond(WContragentGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWContragentGroupByFilterStr(req Req) {
	req.Respond(WContragentGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWContragent(req Req) {
	req.Respond(WContragentGet(req.IntParam))
}

func GetWContragentAll(req Req) {
	req.Respond(WContragentGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWContragentByFilterInt(req Req) {
	req.Respond(WContragentGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWContragentByFilterStr(req Req) {
	req.Respond(WContragentGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWContragentFindByContragentSearchContactSearch(req Req) {
	req.Respond(WContragentFindByContragentSearchContactSearch(req.StrParam))
}

func GetWContact(req Req) {
	req.Respond(WContactGet(req.IntParam))
}

func GetWContactAll(req Req) {
	req.Respond(WContactGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWContactByFilterInt(req Req) {
	req.Respond(WContactGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWContactByFilterStr(req Req) {
	req.Respond(WContactGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWOrderingStatus(req Req) {
	req.Respond(WOrderingStatusGet(req.IntParam))
}

func GetWOrderingStatusAll(req Req) {
	req.Respond(WOrderingStatusGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWOrderingStatusByFilterInt(req Req) {
	req.Respond(WOrderingStatusGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWOrderingStatusByFilterStr(req Req) {
	req.Respond(WOrderingStatusGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWOrdering(req Req) {
	req.Respond(WOrderingGet(req.IntParam))
}

func GetWOrderingAll(req Req) {
	req.Respond(WOrderingGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWOrderingByFilterInt(req Req) {
	req.Respond(WOrderingGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWOrderingByFilterStr(req Req) {
	req.Respond(WOrderingGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWOrderingBetweenCreatedAt(req Req) {
	req.Respond(WOrderingGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWOrderingBetweenDeadlineAt(req Req) {
	req.Respond(WOrderingGetBetweenDeadlineAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWOwner(req Req) {
	req.Respond(WOwnerGet(req.IntParam))
}

func GetWOwnerAll(req Req) {
	req.Respond(WOwnerGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWOwnerByFilterInt(req Req) {
	req.Respond(WOwnerGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWOwnerByFilterStr(req Req) {
	req.Respond(WOwnerGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWInvoice(req Req) {
	req.Respond(WInvoiceGet(req.IntParam))
}

func GetWInvoiceAll(req Req) {
	req.Respond(WInvoiceGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWInvoiceByFilterInt(req Req) {
	req.Respond(WInvoiceGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWInvoiceByFilterStr(req Req) {
	req.Respond(WInvoiceGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWInvoiceBetweenCreatedAt(req Req) {
	req.Respond(WInvoiceGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWItemToInvoice(req Req) {
	req.Respond(WItemToInvoiceGet(req.IntParam))
}

func GetWItemToInvoiceAll(req Req) {
	req.Respond(WItemToInvoiceGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWItemToInvoiceByFilterInt(req Req) {
	req.Respond(WItemToInvoiceGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWItemToInvoiceByFilterStr(req Req) {
	req.Respond(WItemToInvoiceGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWProductToOrderingStatus(req Req) {
	req.Respond(WProductToOrderingStatusGet(req.IntParam))
}

func GetWProductToOrderingStatusAll(req Req) {
	req.Respond(WProductToOrderingStatusGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWProductToOrderingStatusByFilterInt(req Req) {
	req.Respond(WProductToOrderingStatusGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWProductToOrderingStatusByFilterStr(req Req) {
	req.Respond(WProductToOrderingStatusGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWProductToOrdering(req Req) {
	req.Respond(WProductToOrderingGet(req.IntParam))
}

func GetWProductToOrderingAll(req Req) {
	req.Respond(WProductToOrderingGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWProductToOrderingByFilterInt(req Req) {
	req.Respond(WProductToOrderingGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWProductToOrderingByFilterStr(req Req) {
	req.Respond(WProductToOrderingGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWProductToOrderingBetweenUpCreatedAt(req Req) {
	req.Respond(WProductToOrderingGetBetweenUpCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialToOrdering(req Req) {
	req.Respond(WMatherialToOrderingGet(req.IntParam))
}

func GetWMatherialToOrderingAll(req Req) {
	req.Respond(WMatherialToOrderingGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialToOrderingByFilterInt(req Req) {
	req.Respond(WMatherialToOrderingGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialToOrderingByFilterStr(req Req) {
	req.Respond(WMatherialToOrderingGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialToOrderingBetweenUpCreatedAt(req Req) {
	req.Respond(WMatherialToOrderingGetBetweenUpCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialToProduct(req Req) {
	req.Respond(WMatherialToProductGet(req.IntParam))
}

func GetWMatherialToProductAll(req Req) {
	req.Respond(WMatherialToProductGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialToProductByFilterInt(req Req) {
	req.Respond(WMatherialToProductGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialToProductByFilterStr(req Req) {
	req.Respond(WMatherialToProductGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWOperationToOrdering(req Req) {
	req.Respond(WOperationToOrderingGet(req.IntParam))
}

func GetWOperationToOrderingAll(req Req) {
	req.Respond(WOperationToOrderingGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWOperationToOrderingByFilterInt(req Req) {
	req.Respond(WOperationToOrderingGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWOperationToOrderingByFilterStr(req Req) {
	req.Respond(WOperationToOrderingGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWOperationToOrderingBetweenUpCreatedAt(req Req) {
	req.Respond(WOperationToOrderingGetBetweenUpCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWOperationToProduct(req Req) {
	req.Respond(WOperationToProductGet(req.IntParam))
}

func GetWOperationToProductAll(req Req) {
	req.Respond(WOperationToProductGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWOperationToProductByFilterInt(req Req) {
	req.Respond(WOperationToProductGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWOperationToProductByFilterStr(req Req) {
	req.Respond(WOperationToProductGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWProductToProduct(req Req) {
	req.Respond(WProductToProductGet(req.IntParam))
}

func GetWProductToProductAll(req Req) {
	req.Respond(WProductToProductGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWProductToProductByFilterInt(req Req) {
	req.Respond(WProductToProductGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWProductToProductByFilterStr(req Req) {
	req.Respond(WProductToProductGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWCboxCheck(req Req) {
	req.Respond(WCboxCheckGet(req.IntParam))
}

func GetWCboxCheckAll(req Req) {
	req.Respond(WCboxCheckGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWCboxCheckByFilterInt(req Req) {
	req.Respond(WCboxCheckGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWCboxCheckByFilterStr(req Req) {
	req.Respond(WCboxCheckGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWCboxCheckBetweenCreatedAt(req Req) {
	req.Respond(WCboxCheckGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWItemToCboxCheck(req Req) {
	req.Respond(WItemToCboxCheckGet(req.IntParam))
}

func GetWItemToCboxCheckAll(req Req) {
	req.Respond(WItemToCboxCheckGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWItemToCboxCheckByFilterInt(req Req) {
	req.Respond(WItemToCboxCheckGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWItemToCboxCheckByFilterStr(req Req) {
	req.Respond(WItemToCboxCheckGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWCashIn(req Req) {
	req.Respond(WCashInGet(req.IntParam))
}

func GetWCashInAll(req Req) {
	req.Respond(WCashInGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWCashInByFilterInt(req Req) {
	req.Respond(WCashInGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWCashInByFilterStr(req Req) {
	req.Respond(WCashInGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWCashInBetweenCreatedAt(req Req) {
	req.Respond(WCashInGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWCashOut(req Req) {
	req.Respond(WCashOutGet(req.IntParam))
}

func GetWCashOutAll(req Req) {
	req.Respond(WCashOutGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWCashOutByFilterInt(req Req) {
	req.Respond(WCashOutGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWCashOutByFilterStr(req Req) {
	req.Respond(WCashOutGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWCashOutBetweenCreatedAt(req Req) {
	req.Respond(WCashOutGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWWhs(req Req) {
	req.Respond(WWhsGet(req.IntParam))
}

func GetWWhsAll(req Req) {
	req.Respond(WWhsGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWWhsByFilterInt(req Req) {
	req.Respond(WWhsGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWWhsByFilterStr(req Req) {
	req.Respond(WWhsGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWWhsIn(req Req) {
	req.Respond(WWhsInGet(req.IntParam))
}

func GetWWhsInAll(req Req) {
	req.Respond(WWhsInGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWWhsInByFilterInt(req Req) {
	req.Respond(WWhsInGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWWhsInByFilterStr(req Req) {
	req.Respond(WWhsInGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWWhsInBetweenCreatedAt(req Req) {
	req.Respond(WWhsInGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWWhsInBetweenContragentCreatedAt(req Req) {
	req.Respond(WWhsInGetBetweenContragentCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWWhsOut(req Req) {
	req.Respond(WWhsOutGet(req.IntParam))
}

func GetWWhsOutAll(req Req) {
	req.Respond(WWhsOutGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWWhsOutByFilterInt(req Req) {
	req.Respond(WWhsOutGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWWhsOutByFilterStr(req Req) {
	req.Respond(WWhsOutGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWWhsOutBetweenCreatedAt(req Req) {
	req.Respond(WWhsOutGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialToWhsIn(req Req) {
	req.Respond(WMatherialToWhsInGet(req.IntParam))
}

func GetWMatherialToWhsInAll(req Req) {
	req.Respond(WMatherialToWhsInGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialToWhsInByFilterInt(req Req) {
	req.Respond(WMatherialToWhsInGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialToWhsInByFilterStr(req Req) {
	req.Respond(WMatherialToWhsInGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialToWhsOut(req Req) {
	req.Respond(WMatherialToWhsOutGet(req.IntParam))
}

func GetWMatherialToWhsOutAll(req Req) {
	req.Respond(WMatherialToWhsOutGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialToWhsOutByFilterInt(req Req) {
	req.Respond(WMatherialToWhsOutGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialToWhsOutByFilterStr(req Req) {
	req.Respond(WMatherialToWhsOutGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialPart(req Req) {
	req.Respond(WMatherialPartGet(req.IntParam))
}

func GetWMatherialPartAll(req Req) {
	req.Respond(WMatherialPartGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialPartByFilterInt(req Req) {
	req.Respond(WMatherialPartGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialPartByFilterStr(req Req) {
	req.Respond(WMatherialPartGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialPartBetweenCreatedAt(req Req) {
	req.Respond(WMatherialPartGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialPartSlice(req Req) {
	req.Respond(WMatherialPartSliceGet(req.IntParam))
}

func GetWMatherialPartSliceAll(req Req) {
	req.Respond(WMatherialPartSliceGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialPartSliceByFilterInt(req Req) {
	req.Respond(WMatherialPartSliceGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialPartSliceByFilterStr(req Req) {
	req.Respond(WMatherialPartSliceGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWMatherialPartSliceBetweenCreatedAt(req Req) {
	req.Respond(WMatherialPartSliceGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWProjectGroup(req Req) {
	req.Respond(WProjectGroupGet(req.IntParam))
}

func GetWProjectGroupAll(req Req) {
	req.Respond(WProjectGroupGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWProjectGroupByFilterInt(req Req) {
	req.Respond(WProjectGroupGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWProjectGroupByFilterStr(req Req) {
	req.Respond(WProjectGroupGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWProjectStatus(req Req) {
	req.Respond(WProjectStatusGet(req.IntParam))
}

func GetWProjectStatusAll(req Req) {
	req.Respond(WProjectStatusGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWProjectStatusByFilterInt(req Req) {
	req.Respond(WProjectStatusGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWProjectStatusByFilterStr(req Req) {
	req.Respond(WProjectStatusGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWProjectType(req Req) {
	req.Respond(WProjectTypeGet(req.IntParam))
}

func GetWProjectTypeAll(req Req) {
	req.Respond(WProjectTypeGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWProjectTypeByFilterInt(req Req) {
	req.Respond(WProjectTypeGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWProjectTypeByFilterStr(req Req) {
	req.Respond(WProjectTypeGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWProject(req Req) {
	req.Respond(WProjectGet(req.IntParam))
}

func GetWProjectAll(req Req) {
	req.Respond(WProjectGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWProjectByFilterInt(req Req) {
	req.Respond(WProjectGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWProjectByFilterStr(req Req) {
	req.Respond(WProjectGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWProjectFindByProjectInfoContragentNoSearchContactNoSearch(req Req) {
	req.Respond(WProjectFindByProjectInfoContragentNoSearchContactNoSearch(req.StrParam))
}

func GetWProjectBetweenCreatedAt(req Req) {
	req.Respond(WProjectGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWCounter(req Req) {
	req.Respond(WCounterGet(req.IntParam))
}

func GetWCounterAll(req Req) {
	req.Respond(WCounterGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWCounterByFilterInt(req Req) {
	req.Respond(WCounterGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWCounterByFilterStr(req Req) {
	req.Respond(WCounterGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWRecordToCounter(req Req) {
	req.Respond(WRecordToCounterGet(req.IntParam))
}

func GetWRecordToCounterAll(req Req) {
	req.Respond(WRecordToCounterGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWRecordToCounterByFilterInt(req Req) {
	req.Respond(WRecordToCounterGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWRecordToCounterByFilterStr(req Req) {
	req.Respond(WRecordToCounterGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWRecordToCounterBetweenCreatedAt(req Req) {
	req.Respond(WRecordToCounterGetBetweenCreatedAt(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWWmcNumber(req Req) {
	req.Respond(WWmcNumberGet(req.IntParam))
}

func GetWWmcNumberAll(req Req) {
	req.Respond(WWmcNumberGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWWmcNumberByFilterInt(req Req) {
	req.Respond(WWmcNumberGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWWmcNumberByFilterStr(req Req) {
	req.Respond(WWmcNumberGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}

func GetWNumbersToProduct(req Req) {
	req.Respond(WNumbersToProductGet(req.IntParam))
}

func GetWNumbersToProductAll(req Req) {
	req.Respond(WNumbersToProductGetAll(req.WithDeleted, req.DeletedOnly))
}

func GetWNumbersToProductByFilterInt(req Req) {
	req.Respond(WNumbersToProductGetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}

func GetWNumbersToProductByFilterStr(req Req) {
	req.Respond(WNumbersToProductGetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}
