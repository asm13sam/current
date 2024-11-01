package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"asm13sam/tg"

	"github.com/gorilla/mux"
)

func makeRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/document/{id:[0-9]+}",
		WrapAuth(GetDocument, DOC_READ)).Methods("GET")

	r.HandleFunc("/document_get_all",
		WrapAuth(GetDocumentAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/document",
		WrapAuth(CreateDocument, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/document/{id:[0-9]+}",
		WrapAuth(UpdateDocument, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/document/{id:[0-9]+}",
		WrapAuth(DeleteDocument, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/document_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetDocumentByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/document_filter_str/{fs}/{fs2}",
		WrapAuth(GetDocumentByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/measure/{id:[0-9]+}",
		WrapAuth(GetMeasure, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/measure_get_all",
		WrapAuth(GetMeasureAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/measure",
		WrapAuth(CreateMeasure, CATALOG_CREATE)).Methods("POST")

	r.HandleFunc("/measure/{id:[0-9]+}",
		WrapAuth(UpdateMeasure, CATALOG_UPDATE)).Methods("PUT")

	r.HandleFunc("/measure/{id:[0-9]+}",
		WrapAuth(DeleteMeasure, CATALOG_DELETE)).Methods("DELETE")

	r.HandleFunc("/measure_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetMeasureByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/measure_filter_str/{fs}/{fs2}",
		WrapAuth(GetMeasureByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/count_type/{id:[0-9]+}",
		WrapAuth(GetCountType, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/count_type_get_all",
		WrapAuth(GetCountTypeAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/count_type",
		WrapAuth(CreateCountType, CATALOG_CREATE)).Methods("POST")

	r.HandleFunc("/count_type/{id:[0-9]+}",
		WrapAuth(UpdateCountType, CATALOG_UPDATE)).Methods("PUT")

	r.HandleFunc("/count_type/{id:[0-9]+}",
		WrapAuth(DeleteCountType, CATALOG_DELETE)).Methods("DELETE")

	r.HandleFunc("/count_type_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetCountTypeByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/count_type_filter_str/{fs}/{fs2}",
		WrapAuth(GetCountTypeByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/color_group/{id:[0-9]+}",
		WrapAuth(GetColorGroup, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/color_group_get_all",
		WrapAuth(GetColorGroupAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/color_group",
		WrapAuth(CreateColorGroup, CATALOG_CREATE)).Methods("POST")

	r.HandleFunc("/color_group/{id:[0-9]+}",
		WrapAuth(UpdateColorGroup, CATALOG_UPDATE)).Methods("PUT")

	r.HandleFunc("/color_group/{id:[0-9]+}",
		WrapAuth(DeleteColorGroup, CATALOG_DELETE)).Methods("DELETE")

	r.HandleFunc("/color_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetColorGroupByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/color_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetColorGroupByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/color/{id:[0-9]+}",
		WrapAuth(GetColor, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/color_get_all",
		WrapAuth(GetColorAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/color",
		WrapAuth(CreateColor, CATALOG_CREATE)).Methods("POST")

	r.HandleFunc("/color/{id:[0-9]+}",
		WrapAuth(UpdateColor, CATALOG_UPDATE)).Methods("PUT")

	r.HandleFunc("/color/{id:[0-9]+}",
		WrapAuth(DeleteColor, CATALOG_DELETE)).Methods("DELETE")

	r.HandleFunc("/color_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetColorByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/color_filter_str/{fs}/{fs2}",
		WrapAuth(GetColorByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/matherial_group/{id:[0-9]+}",
		WrapAuth(GetMatherialGroup, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/matherial_group_get_all",
		WrapAuth(GetMatherialGroupAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/matherial_group",
		WrapAuth(CreateMatherialGroup, CATALOG_CREATE)).Methods("POST")

	r.HandleFunc("/matherial_group/{id:[0-9]+}",
		WrapAuth(UpdateMatherialGroup, CATALOG_UPDATE)).Methods("PUT")

	r.HandleFunc("/matherial_group/{id:[0-9]+}",
		WrapAuth(DeleteMatherialGroup, CATALOG_DELETE)).Methods("DELETE")

	r.HandleFunc("/matherial_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetMatherialGroupByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/matherial_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetMatherialGroupByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/matherial/{id:[0-9]+}",
		WrapAuth(GetMatherial, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/matherial_get_all",
		WrapAuth(GetMatherialAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/matherial",
		WrapAuth(CreateMatherial, CATALOG_CREATE)).Methods("POST")

	r.HandleFunc("/matherial/{id:[0-9]+}",
		WrapAuth(UpdateMatherial, CATALOG_UPDATE)).Methods("PUT")

	r.HandleFunc("/matherial/{id:[0-9]+}",
		WrapAuth(DeleteMatherial, CATALOG_DELETE)).Methods("DELETE")

	r.HandleFunc("/matherial_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetMatherialByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/matherial_filter_str/{fs}/{fs2}",
		WrapAuth(GetMatherialByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/cash/{id:[0-9]+}",
		WrapAuth(GetCash, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/cash_get_all",
		WrapAuth(GetCashAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/cash",
		WrapAuth(CreateCash, CATALOG_CREATE)).Methods("POST")

	r.HandleFunc("/cash/{id:[0-9]+}",
		WrapAuth(UpdateCash, CATALOG_UPDATE)).Methods("PUT")

	r.HandleFunc("/cash/{id:[0-9]+}",
		WrapAuth(DeleteCash, CATALOG_DELETE)).Methods("DELETE")

	r.HandleFunc("/cash_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetCashByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/cash_filter_str/{fs}/{fs2}",
		WrapAuth(GetCashByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/user_group/{id:[0-9]+}",
		WrapAuth(GetUserGroup, USER_READ)).Methods("GET")

	r.HandleFunc("/user_group_get_all",
		WrapAuth(GetUserGroupAll, USER_READ)).Methods("GET")

	r.HandleFunc("/user_group",
		WrapAuth(CreateUserGroup, USER_CREATE)).Methods("POST")

	r.HandleFunc("/user_group/{id:[0-9]+}",
		WrapAuth(UpdateUserGroup, USER_UPDATE)).Methods("PUT")

	r.HandleFunc("/user_group/{id:[0-9]+}",
		WrapAuth(DeleteUserGroup, USER_DELETE)).Methods("DELETE")

	r.HandleFunc("/user_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetUserGroupByFilterInt, USER_READ)).Methods("GET")

	r.HandleFunc("/user_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetUserGroupByFilterStr, USER_READ)).Methods("GET")

	r.HandleFunc("/user/{id:[0-9]+}",
		WrapAuth(GetUser, USER_READ)).Methods("GET")

	r.HandleFunc("/user_get_all",
		WrapAuth(GetUserAll, USER_READ)).Methods("GET")

	r.HandleFunc("/user",
		WrapAuth(CreateUser, USER_CREATE)).Methods("POST")

	r.HandleFunc("/user/{id:[0-9]+}",
		WrapAuth(UpdateUser, USER_UPDATE)).Methods("PUT")

	r.HandleFunc("/user/{id:[0-9]+}",
		WrapAuth(DeleteUser, USER_DELETE)).Methods("DELETE")

	r.HandleFunc("/user_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetUserByFilterInt, USER_READ)).Methods("GET")

	r.HandleFunc("/user_filter_str/{fs}/{fs2}",
		WrapAuth(GetUserByFilterStr, USER_READ)).Methods("GET")

	r.HandleFunc("/equipment_group/{id:[0-9]+}",
		WrapAuth(GetEquipmentGroup, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/equipment_group_get_all",
		WrapAuth(GetEquipmentGroupAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/equipment_group",
		WrapAuth(CreateEquipmentGroup, CATALOG_CREATE)).Methods("POST")

	r.HandleFunc("/equipment_group/{id:[0-9]+}",
		WrapAuth(UpdateEquipmentGroup, CATALOG_UPDATE)).Methods("PUT")

	r.HandleFunc("/equipment_group/{id:[0-9]+}",
		WrapAuth(DeleteEquipmentGroup, CATALOG_DELETE)).Methods("DELETE")

	r.HandleFunc("/equipment_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetEquipmentGroupByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/equipment_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetEquipmentGroupByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/equipment/{id:[0-9]+}",
		WrapAuth(GetEquipment, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/equipment_get_all",
		WrapAuth(GetEquipmentAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/equipment",
		WrapAuth(CreateEquipment, CATALOG_CREATE)).Methods("POST")

	r.HandleFunc("/equipment/{id:[0-9]+}",
		WrapAuth(UpdateEquipment, CATALOG_UPDATE)).Methods("PUT")

	r.HandleFunc("/equipment/{id:[0-9]+}",
		WrapAuth(DeleteEquipment, CATALOG_DELETE)).Methods("DELETE")

	r.HandleFunc("/equipment_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetEquipmentByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/equipment_filter_str/{fs}/{fs2}",
		WrapAuth(GetEquipmentByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/operation_group/{id:[0-9]+}",
		WrapAuth(GetOperationGroup, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/operation_group_get_all",
		WrapAuth(GetOperationGroupAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/operation_group",
		WrapAuth(CreateOperationGroup, CATALOG_CREATE)).Methods("POST")

	r.HandleFunc("/operation_group/{id:[0-9]+}",
		WrapAuth(UpdateOperationGroup, CATALOG_UPDATE)).Methods("PUT")

	r.HandleFunc("/operation_group/{id:[0-9]+}",
		WrapAuth(DeleteOperationGroup, CATALOG_DELETE)).Methods("DELETE")

	r.HandleFunc("/operation_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetOperationGroupByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/operation_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetOperationGroupByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/operation/{id:[0-9]+}",
		WrapAuth(GetOperation, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/operation_get_all",
		WrapAuth(GetOperationAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/operation",
		WrapAuth(CreateOperation, CATALOG_CREATE)).Methods("POST")

	r.HandleFunc("/operation/{id:[0-9]+}",
		WrapAuth(UpdateOperation, CATALOG_UPDATE)).Methods("PUT")

	r.HandleFunc("/operation/{id:[0-9]+}",
		WrapAuth(DeleteOperation, CATALOG_DELETE)).Methods("DELETE")

	r.HandleFunc("/operation_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetOperationByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/operation_filter_str/{fs}/{fs2}",
		WrapAuth(GetOperationByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/product_group/{id:[0-9]+}",
		WrapAuth(GetProductGroup, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/product_group_get_all",
		WrapAuth(GetProductGroupAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/product_group",
		WrapAuth(CreateProductGroup, CATALOG_CREATE)).Methods("POST")

	r.HandleFunc("/product_group/{id:[0-9]+}",
		WrapAuth(UpdateProductGroup, CATALOG_UPDATE)).Methods("PUT")

	r.HandleFunc("/product_group/{id:[0-9]+}",
		WrapAuth(DeleteProductGroup, CATALOG_DELETE)).Methods("DELETE")

	r.HandleFunc("/product_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetProductGroupByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/product_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetProductGroupByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/product/{id:[0-9]+}",
		WrapAuth(GetProduct, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/product_get_all",
		WrapAuth(GetProductAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/product",
		WrapAuth(CreateProduct, CATALOG_CREATE)).Methods("POST")

	r.HandleFunc("/product/{id:[0-9]+}",
		WrapAuth(UpdateProduct, CATALOG_UPDATE)).Methods("PUT")

	r.HandleFunc("/product/{id:[0-9]+}",
		WrapAuth(DeleteProduct, CATALOG_DELETE)).Methods("DELETE")

	r.HandleFunc("/product_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetProductByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/product_filter_str/{fs}/{fs2}",
		WrapAuth(GetProductByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/contragent_group/{id:[0-9]+}",
		WrapAuth(GetContragentGroup, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/contragent_group_get_all",
		WrapAuth(GetContragentGroupAll, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/contragent_group",
		WrapAuth(CreateContragentGroup, CONTRAGENT_CREATE)).Methods("POST")

	r.HandleFunc("/contragent_group/{id:[0-9]+}",
		WrapAuth(UpdateContragentGroup, CONTRAGENT_UPDATE)).Methods("PUT")

	r.HandleFunc("/contragent_group/{id:[0-9]+}",
		WrapAuth(DeleteContragentGroup, CONTRAGENT_DELETE)).Methods("DELETE")

	r.HandleFunc("/contragent_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetContragentGroupByFilterInt, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/contragent_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetContragentGroupByFilterStr, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/contragent/{id:[0-9]+}",
		WrapAuth(GetContragent, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/contragent_get_all",
		WrapAuth(GetContragentAll, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/contragent",
		WrapAuth(CreateContragent, CONTRAGENT_CREATE)).Methods("POST")

	r.HandleFunc("/contragent/{id:[0-9]+}",
		WrapAuth(UpdateContragent, CONTRAGENT_UPDATE)).Methods("PUT")

	r.HandleFunc("/contragent/{id:[0-9]+}",
		WrapAuth(DeleteContragent, CONTRAGENT_DELETE)).Methods("DELETE")

	r.HandleFunc("/contragent_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetContragentByFilterInt, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/contragent_filter_str/{fs}/{fs2}",
		WrapAuth(GetContragentByFilterStr, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/find_contragent_contragent_search_contact_search/{fs}",
		WrapAuth(GetContragentFindByContragentSearchContactSearch, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/contact/{id:[0-9]+}",
		WrapAuth(GetContact, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/contact_get_all",
		WrapAuth(GetContactAll, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/contact",
		WrapAuth(CreateContact, CONTRAGENT_CREATE)).Methods("POST")

	r.HandleFunc("/contact/{id:[0-9]+}",
		WrapAuth(UpdateContact, CONTRAGENT_UPDATE)).Methods("PUT")

	r.HandleFunc("/contact/{id:[0-9]+}",
		WrapAuth(DeleteContact, CONTRAGENT_DELETE)).Methods("DELETE")

	r.HandleFunc("/contact_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetContactByFilterInt, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/contact_filter_str/{fs}/{fs2}",
		WrapAuth(GetContactByFilterStr, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/ordering_status/{id:[0-9]+}",
		WrapAuth(GetOrderingStatus, DOC_READ)).Methods("GET")

	r.HandleFunc("/ordering_status_get_all",
		WrapAuth(GetOrderingStatusAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/ordering_status",
		WrapAuth(CreateOrderingStatus, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/ordering_status/{id:[0-9]+}",
		WrapAuth(UpdateOrderingStatus, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/ordering_status/{id:[0-9]+}",
		WrapAuth(DeleteOrderingStatus, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/ordering_status_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetOrderingStatusByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/ordering_status_filter_str/{fs}/{fs2}",
		WrapAuth(GetOrderingStatusByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/ordering/{id:[0-9]+}",
		WrapAuth(GetOrdering, DOC_READ)).Methods("GET")

	r.HandleFunc("/ordering_get_all",
		WrapAuth(GetOrderingAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/ordering",
		WrapAuth(CreateOrdering, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/ordering/{id:[0-9]+}",
		WrapAuth(UpdateOrdering, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/ordering/{id:[0-9]+}",
		WrapAuth(DeleteOrdering, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/ordering_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetOrderingByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/ordering_filter_str/{fs}/{fs2}",
		WrapAuth(GetOrderingByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/ordering_between_created_at/{fs}/{fs2}",
		WrapAuth(GetOrderingBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/ordering_between_deadline_at/{fs}/{fs2}",
		WrapAuth(GetOrderingBetweenDeadlineAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/ordering_of_cash_sum_sum_before/{fs}/{id:[0-9]+}/{fs2}",
		WrapAuth(GetOrderingCashSumSumBefore, DOC_READ)).Methods("GET")

	r.HandleFunc("/ordering_sum_filter_by/{fs}/{id:[0-9]+}/{fs2}/{id2:[0-9]+}",
		WrapAuth(GetOrderingSumByFilter, DOC_READ)).Methods("GET")

	r.HandleFunc("/owner/{id:[0-9]+}",
		WrapAuth(GetOwner, OWNER_READ)).Methods("GET")

	r.HandleFunc("/owner_get_all",
		WrapAuth(GetOwnerAll, OWNER_READ)).Methods("GET")

	r.HandleFunc("/owner",
		WrapAuth(CreateOwner, OWNER_CREATE)).Methods("POST")

	r.HandleFunc("/owner/{id:[0-9]+}",
		WrapAuth(UpdateOwner, OWNER_UPDATE)).Methods("PUT")

	r.HandleFunc("/owner/{id:[0-9]+}",
		WrapAuth(DeleteOwner, OWNER_DELETE)).Methods("DELETE")

	r.HandleFunc("/owner_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetOwnerByFilterInt, OWNER_READ)).Methods("GET")

	r.HandleFunc("/owner_filter_str/{fs}/{fs2}",
		WrapAuth(GetOwnerByFilterStr, OWNER_READ)).Methods("GET")

	r.HandleFunc("/invoice/{id:[0-9]+}",
		WrapAuth(GetInvoice, DOC_READ)).Methods("GET")

	r.HandleFunc("/invoice_get_all",
		WrapAuth(GetInvoiceAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/invoice",
		WrapAuth(CreateInvoice, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/invoice/{id:[0-9]+}",
		WrapAuth(UpdateInvoice, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/invoice/{id:[0-9]+}",
		WrapAuth(DeleteInvoice, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/invoice_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetInvoiceByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/invoice_filter_str/{fs}/{fs2}",
		WrapAuth(GetInvoiceByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/invoice_between_created_at/{fs}/{fs2}",
		WrapAuth(GetInvoiceBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/invoice_of_cash_sum_sum_before/{fs}/{id:[0-9]+}/{fs2}",
		WrapAuth(GetInvoiceCashSumSumBefore, DOC_READ)).Methods("GET")

	r.HandleFunc("/invoice_sum_filter_by/{fs}/{id:[0-9]+}/{fs2}/{id2:[0-9]+}",
		WrapAuth(GetInvoiceSumByFilter, DOC_READ)).Methods("GET")

	r.HandleFunc("/item_to_invoice/{id:[0-9]+}",
		WrapAuth(GetItemToInvoice, DOC_READ)).Methods("GET")

	r.HandleFunc("/item_to_invoice_get_all",
		WrapAuth(GetItemToInvoiceAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/item_to_invoice",
		WrapAuth(CreateItemToInvoice, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/item_to_invoice/{id:[0-9]+}",
		WrapAuth(UpdateItemToInvoice, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/item_to_invoice/{id:[0-9]+}",
		WrapAuth(DeleteItemToInvoice, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/item_to_invoice_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetItemToInvoiceByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/item_to_invoice_filter_str/{fs}/{fs2}",
		WrapAuth(GetItemToInvoiceByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/realized/item_to_invoice/{id:[0-9]+}",
		WrapAuth(RealizedItemToInvoice, DOC_CREATE)).Methods("GET")

	r.HandleFunc("/item_to_invoice_of_cost_sum_before/{fs}/{id:[0-9]+}/{fs2}",
		WrapAuth(GetItemToInvoiceCostSumBefore, DOC_READ)).Methods("GET")

	r.HandleFunc("/item_to_invoice_sum_filter_by/{fs}/{id:[0-9]+}/{fs2}/{id2:[0-9]+}",
		WrapAuth(GetItemToInvoiceSumByFilter, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_ordering_status/{id:[0-9]+}",
		WrapAuth(GetProductToOrderingStatus, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_ordering_status_get_all",
		WrapAuth(GetProductToOrderingStatusAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_ordering_status",
		WrapAuth(CreateProductToOrderingStatus, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/product_to_ordering_status/{id:[0-9]+}",
		WrapAuth(UpdateProductToOrderingStatus, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/product_to_ordering_status/{id:[0-9]+}",
		WrapAuth(DeleteProductToOrderingStatus, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/product_to_ordering_status_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetProductToOrderingStatusByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_ordering_status_filter_str/{fs}/{fs2}",
		WrapAuth(GetProductToOrderingStatusByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_ordering/{id:[0-9]+}",
		WrapAuth(GetProductToOrdering, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_ordering_get_all",
		WrapAuth(GetProductToOrderingAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_ordering",
		WrapAuth(CreateProductToOrdering, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/product_to_ordering/{id:[0-9]+}",
		WrapAuth(UpdateProductToOrdering, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/product_to_ordering/{id:[0-9]+}",
		WrapAuth(DeleteProductToOrdering, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/product_to_ordering_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetProductToOrderingByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_ordering_filter_str/{fs}/{fs2}",
		WrapAuth(GetProductToOrderingByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_ordering_of_cost_sum_before/{fs}/{id:[0-9]+}/{fs2}",
		WrapAuth(GetProductToOrderingCostSumBefore, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_ordering_sum_filter_by/{fs}/{id:[0-9]+}/{fs2}/{id2:[0-9]+}",
		WrapAuth(GetProductToOrderingSumByFilter, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_ordering/{id:[0-9]+}",
		WrapAuth(GetMatherialToOrdering, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_ordering_get_all",
		WrapAuth(GetMatherialToOrderingAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_ordering",
		WrapAuth(CreateMatherialToOrdering, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/matherial_to_ordering/{id:[0-9]+}",
		WrapAuth(UpdateMatherialToOrdering, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/matherial_to_ordering/{id:[0-9]+}",
		WrapAuth(DeleteMatherialToOrdering, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/matherial_to_ordering_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetMatherialToOrderingByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_ordering_filter_str/{fs}/{fs2}",
		WrapAuth(GetMatherialToOrderingByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_ordering_of_cost_sum_before/{fs}/{id:[0-9]+}/{fs2}",
		WrapAuth(GetMatherialToOrderingCostSumBefore, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_ordering_sum_filter_by/{fs}/{id:[0-9]+}/{fs2}/{id2:[0-9]+}",
		WrapAuth(GetMatherialToOrderingSumByFilter, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_product/{id:[0-9]+}",
		WrapAuth(GetMatherialToProduct, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_product_get_all",
		WrapAuth(GetMatherialToProductAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_product",
		WrapAuth(CreateMatherialToProduct, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/matherial_to_product/{id:[0-9]+}",
		WrapAuth(UpdateMatherialToProduct, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/matherial_to_product/{id:[0-9]+}",
		WrapAuth(DeleteMatherialToProduct, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/matherial_to_product_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetMatherialToProductByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_product_filter_str/{fs}/{fs2}",
		WrapAuth(GetMatherialToProductByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_product_of_cost_sum_before/{fs}/{id:[0-9]+}/{fs2}",
		WrapAuth(GetMatherialToProductCostSumBefore, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_product_sum_filter_by/{fs}/{id:[0-9]+}/{fs2}/{id2:[0-9]+}",
		WrapAuth(GetMatherialToProductSumByFilter, DOC_READ)).Methods("GET")

	r.HandleFunc("/operation_to_ordering/{id:[0-9]+}",
		WrapAuth(GetOperationToOrdering, DOC_READ)).Methods("GET")

	r.HandleFunc("/operation_to_ordering_get_all",
		WrapAuth(GetOperationToOrderingAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/operation_to_ordering",
		WrapAuth(CreateOperationToOrdering, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/operation_to_ordering/{id:[0-9]+}",
		WrapAuth(UpdateOperationToOrdering, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/operation_to_ordering/{id:[0-9]+}",
		WrapAuth(DeleteOperationToOrdering, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/operation_to_ordering_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetOperationToOrderingByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/operation_to_ordering_filter_str/{fs}/{fs2}",
		WrapAuth(GetOperationToOrderingByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/operation_to_ordering_of_cost_sum_before/{fs}/{id:[0-9]+}/{fs2}",
		WrapAuth(GetOperationToOrderingCostSumBefore, DOC_READ)).Methods("GET")

	r.HandleFunc("/operation_to_ordering_sum_filter_by/{fs}/{id:[0-9]+}/{fs2}/{id2:[0-9]+}",
		WrapAuth(GetOperationToOrderingSumByFilter, DOC_READ)).Methods("GET")

	r.HandleFunc("/operation_to_product/{id:[0-9]+}",
		WrapAuth(GetOperationToProduct, DOC_READ)).Methods("GET")

	r.HandleFunc("/operation_to_product_get_all",
		WrapAuth(GetOperationToProductAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/operation_to_product",
		WrapAuth(CreateOperationToProduct, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/operation_to_product/{id:[0-9]+}",
		WrapAuth(UpdateOperationToProduct, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/operation_to_product/{id:[0-9]+}",
		WrapAuth(DeleteOperationToProduct, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/operation_to_product_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetOperationToProductByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/operation_to_product_filter_str/{fs}/{fs2}",
		WrapAuth(GetOperationToProductByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/operation_to_product_of_cost_sum_before/{fs}/{id:[0-9]+}/{fs2}",
		WrapAuth(GetOperationToProductCostSumBefore, DOC_READ)).Methods("GET")

	r.HandleFunc("/operation_to_product_sum_filter_by/{fs}/{id:[0-9]+}/{fs2}/{id2:[0-9]+}",
		WrapAuth(GetOperationToProductSumByFilter, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_product/{id:[0-9]+}",
		WrapAuth(GetProductToProduct, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_product_get_all",
		WrapAuth(GetProductToProductAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_product",
		WrapAuth(CreateProductToProduct, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/product_to_product/{id:[0-9]+}",
		WrapAuth(UpdateProductToProduct, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/product_to_product/{id:[0-9]+}",
		WrapAuth(DeleteProductToProduct, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/product_to_product_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetProductToProductByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_product_filter_str/{fs}/{fs2}",
		WrapAuth(GetProductToProductByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_product_of_cost_sum_before/{fs}/{id:[0-9]+}/{fs2}",
		WrapAuth(GetProductToProductCostSumBefore, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_to_product_sum_filter_by/{fs}/{id:[0-9]+}/{fs2}/{id2:[0-9]+}",
		WrapAuth(GetProductToProductSumByFilter, DOC_READ)).Methods("GET")

	r.HandleFunc("/cbox_check/{id:[0-9]+}",
		WrapAuth(GetCboxCheck, DOC_READ)).Methods("GET")

	r.HandleFunc("/cbox_check_get_all",
		WrapAuth(GetCboxCheckAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/cbox_check",
		WrapAuth(CreateCboxCheck, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/cbox_check/{id:[0-9]+}",
		WrapAuth(UpdateCboxCheck, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/cbox_check/{id:[0-9]+}",
		WrapAuth(DeleteCboxCheck, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/cbox_check_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetCboxCheckByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/cbox_check_filter_str/{fs}/{fs2}",
		WrapAuth(GetCboxCheckByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/cbox_check_between_created_at/{fs}/{fs2}",
		WrapAuth(GetCboxCheckBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/item_to_cbox_check/{id:[0-9]+}",
		WrapAuth(GetItemToCboxCheck, DOC_READ)).Methods("GET")

	r.HandleFunc("/item_to_cbox_check_get_all",
		WrapAuth(GetItemToCboxCheckAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/item_to_cbox_check",
		WrapAuth(CreateItemToCboxCheck, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/item_to_cbox_check/{id:[0-9]+}",
		WrapAuth(UpdateItemToCboxCheck, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/item_to_cbox_check/{id:[0-9]+}",
		WrapAuth(DeleteItemToCboxCheck, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/item_to_cbox_check_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetItemToCboxCheckByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/item_to_cbox_check_filter_str/{fs}/{fs2}",
		WrapAuth(GetItemToCboxCheckByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/item_to_cbox_check_of_cost_sum_before/{fs}/{id:[0-9]+}/{fs2}",
		WrapAuth(GetItemToCboxCheckCostSumBefore, DOC_READ)).Methods("GET")

	r.HandleFunc("/item_to_cbox_check_sum_filter_by/{fs}/{id:[0-9]+}/{fs2}/{id2:[0-9]+}",
		WrapAuth(GetItemToCboxCheckSumByFilter, DOC_READ)).Methods("GET")

	r.HandleFunc("/cash_in/{id:[0-9]+}",
		WrapAuth(GetCashIn, DOC_READ)).Methods("GET")

	r.HandleFunc("/cash_in_get_all",
		WrapAuth(GetCashInAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/cash_in",
		WrapAuth(CreateCashIn, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/cash_in/{id:[0-9]+}",
		WrapAuth(UpdateCashIn, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/cash_in/{id:[0-9]+}",
		WrapAuth(DeleteCashIn, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/cash_in_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetCashInByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/cash_in_filter_str/{fs}/{fs2}",
		WrapAuth(GetCashInByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/cash_in_between_created_at/{fs}/{fs2}",
		WrapAuth(GetCashInBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/cash_in_of_cash_sum_sum_before/{fs}/{id:[0-9]+}/{fs2}",
		WrapAuth(GetCashInCashSumSumBefore, DOC_READ)).Methods("GET")

	r.HandleFunc("/cash_in_sum_filter_by/{fs}/{id:[0-9]+}/{fs2}/{id2:[0-9]+}",
		WrapAuth(GetCashInSumByFilter, DOC_READ)).Methods("GET")

	r.HandleFunc("/cash_out/{id:[0-9]+}",
		WrapAuth(GetCashOut, DOC_READ)).Methods("GET")

	r.HandleFunc("/cash_out_get_all",
		WrapAuth(GetCashOutAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/cash_out",
		WrapAuth(CreateCashOut, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/cash_out/{id:[0-9]+}",
		WrapAuth(UpdateCashOut, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/cash_out/{id:[0-9]+}",
		WrapAuth(DeleteCashOut, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/cash_out_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetCashOutByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/cash_out_filter_str/{fs}/{fs2}",
		WrapAuth(GetCashOutByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/cash_out_between_created_at/{fs}/{fs2}",
		WrapAuth(GetCashOutBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/cash_out_of_cash_sum_sum_before/{fs}/{id:[0-9]+}/{fs2}",
		WrapAuth(GetCashOutCashSumSumBefore, DOC_READ)).Methods("GET")

	r.HandleFunc("/cash_out_sum_filter_by/{fs}/{id:[0-9]+}/{fs2}/{id2:[0-9]+}",
		WrapAuth(GetCashOutSumByFilter, DOC_READ)).Methods("GET")

	r.HandleFunc("/whs/{id:[0-9]+}",
		WrapAuth(GetWhs, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/whs_get_all",
		WrapAuth(GetWhsAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/whs",
		WrapAuth(CreateWhs, CATALOG_CREATE)).Methods("POST")

	r.HandleFunc("/whs/{id:[0-9]+}",
		WrapAuth(UpdateWhs, CATALOG_UPDATE)).Methods("PUT")

	r.HandleFunc("/whs/{id:[0-9]+}",
		WrapAuth(DeleteWhs, CATALOG_DELETE)).Methods("DELETE")

	r.HandleFunc("/whs_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWhsByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/whs_filter_str/{fs}/{fs2}",
		WrapAuth(GetWhsByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/whs_in/{id:[0-9]+}",
		WrapAuth(GetWhsIn, DOC_READ)).Methods("GET")

	r.HandleFunc("/whs_in_get_all",
		WrapAuth(GetWhsInAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/whs_in",
		WrapAuth(CreateWhsIn, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/whs_in/{id:[0-9]+}",
		WrapAuth(UpdateWhsIn, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/whs_in/{id:[0-9]+}",
		WrapAuth(DeleteWhsIn, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/whs_in_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWhsInByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/whs_in_filter_str/{fs}/{fs2}",
		WrapAuth(GetWhsInByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/whs_in_between_created_at/{fs}/{fs2}",
		WrapAuth(GetWhsInBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/whs_in_between_contragent_created_at/{fs}/{fs2}",
		WrapAuth(GetWhsInBetweenContragentCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/whs_in_of_whs_sum_sum_before/{fs}/{id:[0-9]+}/{fs2}",
		WrapAuth(GetWhsInWhsSumSumBefore, DOC_READ)).Methods("GET")

	r.HandleFunc("/whs_in_sum_filter_by/{fs}/{id:[0-9]+}/{fs2}/{id2:[0-9]+}",
		WrapAuth(GetWhsInSumByFilter, DOC_READ)).Methods("GET")

	r.HandleFunc("/whs_out/{id:[0-9]+}",
		WrapAuth(GetWhsOut, DOC_READ)).Methods("GET")

	r.HandleFunc("/whs_out_get_all",
		WrapAuth(GetWhsOutAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/whs_out",
		WrapAuth(CreateWhsOut, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/whs_out/{id:[0-9]+}",
		WrapAuth(UpdateWhsOut, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/whs_out/{id:[0-9]+}",
		WrapAuth(DeleteWhsOut, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/whs_out_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWhsOutByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/whs_out_filter_str/{fs}/{fs2}",
		WrapAuth(GetWhsOutByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/whs_out_between_created_at/{fs}/{fs2}",
		WrapAuth(GetWhsOutBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/whs_out_of_whs_sum_sum_before/{fs}/{id:[0-9]+}/{fs2}",
		WrapAuth(GetWhsOutWhsSumSumBefore, DOC_READ)).Methods("GET")

	r.HandleFunc("/whs_out_sum_filter_by/{fs}/{id:[0-9]+}/{fs2}/{id2:[0-9]+}",
		WrapAuth(GetWhsOutSumByFilter, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_whs_in/{id:[0-9]+}",
		WrapAuth(GetMatherialToWhsIn, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_whs_in_get_all",
		WrapAuth(GetMatherialToWhsInAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_whs_in",
		WrapAuth(CreateMatherialToWhsIn, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/matherial_to_whs_in/{id:[0-9]+}",
		WrapAuth(UpdateMatherialToWhsIn, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/matherial_to_whs_in/{id:[0-9]+}",
		WrapAuth(DeleteMatherialToWhsIn, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/matherial_to_whs_in_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetMatherialToWhsInByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_whs_in_filter_str/{fs}/{fs2}",
		WrapAuth(GetMatherialToWhsInByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/realized/matherial_to_whs_in/{id:[0-9]+}",
		WrapAuth(RealizedMatherialToWhsIn, DOC_CREATE)).Methods("GET")

	r.HandleFunc("/matherial_to_whs_out/{id:[0-9]+}",
		WrapAuth(GetMatherialToWhsOut, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_whs_out_get_all",
		WrapAuth(GetMatherialToWhsOutAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_whs_out",
		WrapAuth(CreateMatherialToWhsOut, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/matherial_to_whs_out/{id:[0-9]+}",
		WrapAuth(UpdateMatherialToWhsOut, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/matherial_to_whs_out/{id:[0-9]+}",
		WrapAuth(DeleteMatherialToWhsOut, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/matherial_to_whs_out_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetMatherialToWhsOutByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_to_whs_out_filter_str/{fs}/{fs2}",
		WrapAuth(GetMatherialToWhsOutByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/realized/matherial_to_whs_out/{id:[0-9]+}",
		WrapAuth(RealizedMatherialToWhsOut, DOC_CREATE)).Methods("GET")

	r.HandleFunc("/matherial_part/{id:[0-9]+}",
		WrapAuth(GetMatherialPart, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_part_get_all",
		WrapAuth(GetMatherialPartAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_part",
		WrapAuth(CreateMatherialPart, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/matherial_part/{id:[0-9]+}",
		WrapAuth(UpdateMatherialPart, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/matherial_part/{id:[0-9]+}",
		WrapAuth(DeleteMatherialPart, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/matherial_part_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetMatherialPartByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_part_filter_str/{fs}/{fs2}",
		WrapAuth(GetMatherialPartByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_part_between_created_at/{fs}/{fs2}",
		WrapAuth(GetMatherialPartBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_part_slice/{id:[0-9]+}",
		WrapAuth(GetMatherialPartSlice, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_part_slice_get_all",
		WrapAuth(GetMatherialPartSliceAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_part_slice",
		WrapAuth(CreateMatherialPartSlice, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/matherial_part_slice/{id:[0-9]+}",
		WrapAuth(UpdateMatherialPartSlice, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/matherial_part_slice/{id:[0-9]+}",
		WrapAuth(DeleteMatherialPartSlice, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/matherial_part_slice_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetMatherialPartSliceByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_part_slice_filter_str/{fs}/{fs2}",
		WrapAuth(GetMatherialPartSliceByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/matherial_part_slice_between_created_at/{fs}/{fs2}",
		WrapAuth(GetMatherialPartSliceBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_group/{id:[0-9]+}",
		WrapAuth(GetProjectGroup, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_group_get_all",
		WrapAuth(GetProjectGroupAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_group",
		WrapAuth(CreateProjectGroup, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/project_group/{id:[0-9]+}",
		WrapAuth(UpdateProjectGroup, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/project_group/{id:[0-9]+}",
		WrapAuth(DeleteProjectGroup, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/project_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetProjectGroupByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetProjectGroupByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_status/{id:[0-9]+}",
		WrapAuth(GetProjectStatus, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_status_get_all",
		WrapAuth(GetProjectStatusAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_status",
		WrapAuth(CreateProjectStatus, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/project_status/{id:[0-9]+}",
		WrapAuth(UpdateProjectStatus, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/project_status/{id:[0-9]+}",
		WrapAuth(DeleteProjectStatus, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/project_status_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetProjectStatusByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_status_filter_str/{fs}/{fs2}",
		WrapAuth(GetProjectStatusByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_type/{id:[0-9]+}",
		WrapAuth(GetProjectType, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_type_get_all",
		WrapAuth(GetProjectTypeAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_type",
		WrapAuth(CreateProjectType, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/project_type/{id:[0-9]+}",
		WrapAuth(UpdateProjectType, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/project_type/{id:[0-9]+}",
		WrapAuth(DeleteProjectType, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/project_type_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetProjectTypeByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_type_filter_str/{fs}/{fs2}",
		WrapAuth(GetProjectTypeByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/project/{id:[0-9]+}",
		WrapAuth(GetProject, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_get_all",
		WrapAuth(GetProjectAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/project",
		WrapAuth(CreateProject, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/project/{id:[0-9]+}",
		WrapAuth(UpdateProject, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/project/{id:[0-9]+}",
		WrapAuth(DeleteProject, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/project_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetProjectByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_filter_str/{fs}/{fs2}",
		WrapAuth(GetProjectByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/find_project_project_info_contragent_no_search_contact_no_search/{fs}",
		WrapAuth(GetProjectFindByProjectInfoContragentNoSearchContactNoSearch, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_between_created_at/{fs}/{fs2}",
		WrapAuth(GetProjectBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/counter/{id:[0-9]+}",
		WrapAuth(GetCounter, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/counter_get_all",
		WrapAuth(GetCounterAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/counter",
		WrapAuth(CreateCounter, CATALOG_CREATE)).Methods("POST")

	r.HandleFunc("/counter/{id:[0-9]+}",
		WrapAuth(UpdateCounter, CATALOG_UPDATE)).Methods("PUT")

	r.HandleFunc("/counter/{id:[0-9]+}",
		WrapAuth(DeleteCounter, CATALOG_DELETE)).Methods("DELETE")

	r.HandleFunc("/counter_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetCounterByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/counter_filter_str/{fs}/{fs2}",
		WrapAuth(GetCounterByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/record_to_counter/{id:[0-9]+}",
		WrapAuth(GetRecordToCounter, DOC_READ)).Methods("GET")

	r.HandleFunc("/record_to_counter_get_all",
		WrapAuth(GetRecordToCounterAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/record_to_counter",
		WrapAuth(CreateRecordToCounter, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/record_to_counter/{id:[0-9]+}",
		WrapAuth(UpdateRecordToCounter, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/record_to_counter_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetRecordToCounterByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/record_to_counter_filter_str/{fs}/{fs2}",
		WrapAuth(GetRecordToCounterByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/record_to_counter_between_created_at/{fs}/{fs2}",
		WrapAuth(GetRecordToCounterBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/record_to_counter_of_number_sum_before/{fs}/{id:[0-9]+}/{fs2}",
		WrapAuth(GetRecordToCounterNumberSumBefore, DOC_READ)).Methods("GET")

	r.HandleFunc("/record_to_counter_sum_filter_by/{fs}/{id:[0-9]+}/{fs2}/{id2:[0-9]+}",
		WrapAuth(GetRecordToCounterSumByFilter, DOC_READ)).Methods("GET")

	r.HandleFunc("/wmc_number/{id:[0-9]+}",
		WrapAuth(GetWmcNumber, DOC_READ)).Methods("GET")

	r.HandleFunc("/wmc_number_get_all",
		WrapAuth(GetWmcNumberAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/wmc_number",
		WrapAuth(CreateWmcNumber, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/wmc_number/{id:[0-9]+}",
		WrapAuth(UpdateWmcNumber, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/wmc_number/{id:[0-9]+}",
		WrapAuth(DeleteWmcNumber, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/wmc_number_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWmcNumberByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/wmc_number_filter_str/{fs}/{fs2}",
		WrapAuth(GetWmcNumberByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/numbers_to_product/{id:[0-9]+}",
		WrapAuth(GetNumbersToProduct, DOC_READ)).Methods("GET")

	r.HandleFunc("/numbers_to_product_get_all",
		WrapAuth(GetNumbersToProductAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/numbers_to_product",
		WrapAuth(CreateNumbersToProduct, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/numbers_to_product/{id:[0-9]+}",
		WrapAuth(UpdateNumbersToProduct, DOC_UPDATE)).Methods("PUT")

	r.HandleFunc("/numbers_to_product/{id:[0-9]+}",
		WrapAuth(DeleteNumbersToProduct, DOC_DELETE)).Methods("DELETE")

	r.HandleFunc("/numbers_to_product_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetNumbersToProductByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/numbers_to_product_filter_str/{fs}/{fs2}",
		WrapAuth(GetNumbersToProductByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_document/{id:[0-9]+}",
		WrapAuth(GetWDocument, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_document_get_all",
		WrapAuth(GetWDocumentAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_document_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWDocumentByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_document_filter_str/{fs}/{fs2}",
		WrapAuth(GetWDocumentByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_measure/{id:[0-9]+}",
		WrapAuth(GetWMeasure, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_measure_get_all",
		WrapAuth(GetWMeasureAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_measure_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWMeasureByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_measure_filter_str/{fs}/{fs2}",
		WrapAuth(GetWMeasureByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_count_type/{id:[0-9]+}",
		WrapAuth(GetWCountType, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_count_type_get_all",
		WrapAuth(GetWCountTypeAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_count_type_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWCountTypeByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_count_type_filter_str/{fs}/{fs2}",
		WrapAuth(GetWCountTypeByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_color_group/{id:[0-9]+}",
		WrapAuth(GetWColorGroup, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_color_group_get_all",
		WrapAuth(GetWColorGroupAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_color_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWColorGroupByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_color_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetWColorGroupByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_color/{id:[0-9]+}",
		WrapAuth(GetWColor, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_color_get_all",
		WrapAuth(GetWColorAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_color_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWColorByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_color_filter_str/{fs}/{fs2}",
		WrapAuth(GetWColorByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_group/{id:[0-9]+}",
		WrapAuth(GetWMatherialGroup, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_group_get_all",
		WrapAuth(GetWMatherialGroupAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWMatherialGroupByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetWMatherialGroupByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_matherial/{id:[0-9]+}",
		WrapAuth(GetWMatherial, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_get_all",
		WrapAuth(GetWMatherialAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWMatherialByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_filter_str/{fs}/{fs2}",
		WrapAuth(GetWMatherialByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_cash/{id:[0-9]+}",
		WrapAuth(GetWCash, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_cash_get_all",
		WrapAuth(GetWCashAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_cash_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWCashByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_cash_filter_str/{fs}/{fs2}",
		WrapAuth(GetWCashByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_user_group/{id:[0-9]+}",
		WrapAuth(GetWUserGroup, USER_READ)).Methods("GET")

	r.HandleFunc("/w_user_group_get_all",
		WrapAuth(GetWUserGroupAll, USER_READ)).Methods("GET")

	r.HandleFunc("/w_user_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWUserGroupByFilterInt, USER_READ)).Methods("GET")

	r.HandleFunc("/w_user_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetWUserGroupByFilterStr, USER_READ)).Methods("GET")

	r.HandleFunc("/w_user/{id:[0-9]+}",
		WrapAuth(GetWUser, USER_READ)).Methods("GET")

	r.HandleFunc("/w_user_get_all",
		WrapAuth(GetWUserAll, USER_READ)).Methods("GET")

	r.HandleFunc("/w_user_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWUserByFilterInt, USER_READ)).Methods("GET")

	r.HandleFunc("/w_user_filter_str/{fs}/{fs2}",
		WrapAuth(GetWUserByFilterStr, USER_READ)).Methods("GET")

	r.HandleFunc("/w_equipment_group/{id:[0-9]+}",
		WrapAuth(GetWEquipmentGroup, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_equipment_group_get_all",
		WrapAuth(GetWEquipmentGroupAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_equipment_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWEquipmentGroupByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_equipment_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetWEquipmentGroupByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_equipment/{id:[0-9]+}",
		WrapAuth(GetWEquipment, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_equipment_get_all",
		WrapAuth(GetWEquipmentAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_equipment_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWEquipmentByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_equipment_filter_str/{fs}/{fs2}",
		WrapAuth(GetWEquipmentByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_operation_group/{id:[0-9]+}",
		WrapAuth(GetWOperationGroup, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_operation_group_get_all",
		WrapAuth(GetWOperationGroupAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_operation_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWOperationGroupByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_operation_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetWOperationGroupByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_operation/{id:[0-9]+}",
		WrapAuth(GetWOperation, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_operation_get_all",
		WrapAuth(GetWOperationAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_operation_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWOperationByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_operation_filter_str/{fs}/{fs2}",
		WrapAuth(GetWOperationByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_product_group/{id:[0-9]+}",
		WrapAuth(GetWProductGroup, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_product_group_get_all",
		WrapAuth(GetWProductGroupAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_product_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWProductGroupByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_product_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetWProductGroupByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_product/{id:[0-9]+}",
		WrapAuth(GetWProduct, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_product_get_all",
		WrapAuth(GetWProductAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_product_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWProductByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_product_filter_str/{fs}/{fs2}",
		WrapAuth(GetWProductByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_contragent_group/{id:[0-9]+}",
		WrapAuth(GetWContragentGroup, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/w_contragent_group_get_all",
		WrapAuth(GetWContragentGroupAll, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/w_contragent_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWContragentGroupByFilterInt, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/w_contragent_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetWContragentGroupByFilterStr, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/w_contragent/{id:[0-9]+}",
		WrapAuth(GetWContragent, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/w_contragent_get_all",
		WrapAuth(GetWContragentAll, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/w_contragent_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWContragentByFilterInt, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/w_contragent_filter_str/{fs}/{fs2}",
		WrapAuth(GetWContragentByFilterStr, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/w_find_contragent_contragent_search_contact_search/{fs}",
		WrapAuth(GetWContragentFindByContragentSearchContactSearch, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/w_contact/{id:[0-9]+}",
		WrapAuth(GetWContact, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/w_contact_get_all",
		WrapAuth(GetWContactAll, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/w_contact_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWContactByFilterInt, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/w_contact_filter_str/{fs}/{fs2}",
		WrapAuth(GetWContactByFilterStr, CONTRAGENT_READ)).Methods("GET")

	r.HandleFunc("/w_ordering_status/{id:[0-9]+}",
		WrapAuth(GetWOrderingStatus, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_ordering_status_get_all",
		WrapAuth(GetWOrderingStatusAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_ordering_status_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWOrderingStatusByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_ordering_status_filter_str/{fs}/{fs2}",
		WrapAuth(GetWOrderingStatusByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_ordering/{id:[0-9]+}",
		WrapAuth(GetWOrdering, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_ordering_get_all",
		WrapAuth(GetWOrderingAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_ordering_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWOrderingByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_ordering_filter_str/{fs}/{fs2}",
		WrapAuth(GetWOrderingByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_ordering_between_created_at/{fs}/{fs2}",
		WrapAuth(GetWOrderingBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_ordering_between_deadline_at/{fs}/{fs2}",
		WrapAuth(GetWOrderingBetweenDeadlineAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_owner/{id:[0-9]+}",
		WrapAuth(GetWOwner, OWNER_READ)).Methods("GET")

	r.HandleFunc("/w_owner_get_all",
		WrapAuth(GetWOwnerAll, OWNER_READ)).Methods("GET")

	r.HandleFunc("/w_owner_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWOwnerByFilterInt, OWNER_READ)).Methods("GET")

	r.HandleFunc("/w_owner_filter_str/{fs}/{fs2}",
		WrapAuth(GetWOwnerByFilterStr, OWNER_READ)).Methods("GET")

	r.HandleFunc("/w_invoice/{id:[0-9]+}",
		WrapAuth(GetWInvoice, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_invoice_get_all",
		WrapAuth(GetWInvoiceAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_invoice_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWInvoiceByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_invoice_filter_str/{fs}/{fs2}",
		WrapAuth(GetWInvoiceByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_invoice_between_created_at/{fs}/{fs2}",
		WrapAuth(GetWInvoiceBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_item_to_invoice/{id:[0-9]+}",
		WrapAuth(GetWItemToInvoice, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_item_to_invoice_get_all",
		WrapAuth(GetWItemToInvoiceAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_item_to_invoice_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWItemToInvoiceByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_item_to_invoice_filter_str/{fs}/{fs2}",
		WrapAuth(GetWItemToInvoiceByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_product_to_ordering_status/{id:[0-9]+}",
		WrapAuth(GetWProductToOrderingStatus, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_product_to_ordering_status_get_all",
		WrapAuth(GetWProductToOrderingStatusAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_product_to_ordering_status_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWProductToOrderingStatusByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_product_to_ordering_status_filter_str/{fs}/{fs2}",
		WrapAuth(GetWProductToOrderingStatusByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_product_to_ordering/{id:[0-9]+}",
		WrapAuth(GetWProductToOrdering, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_product_to_ordering_get_all",
		WrapAuth(GetWProductToOrderingAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_product_to_ordering_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWProductToOrderingByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_product_to_ordering_filter_str/{fs}/{fs2}",
		WrapAuth(GetWProductToOrderingByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_product_to_ordering_between_up_created_at/{fs}/{fs2}",
		WrapAuth(GetWProductToOrderingBetweenUpCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_ordering/{id:[0-9]+}",
		WrapAuth(GetWMatherialToOrdering, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_ordering_get_all",
		WrapAuth(GetWMatherialToOrderingAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_ordering_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWMatherialToOrderingByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_ordering_filter_str/{fs}/{fs2}",
		WrapAuth(GetWMatherialToOrderingByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_ordering_between_up_created_at/{fs}/{fs2}",
		WrapAuth(GetWMatherialToOrderingBetweenUpCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_product/{id:[0-9]+}",
		WrapAuth(GetWMatherialToProduct, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_product_get_all",
		WrapAuth(GetWMatherialToProductAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_product_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWMatherialToProductByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_product_filter_str/{fs}/{fs2}",
		WrapAuth(GetWMatherialToProductByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_operation_to_ordering/{id:[0-9]+}",
		WrapAuth(GetWOperationToOrdering, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_operation_to_ordering_get_all",
		WrapAuth(GetWOperationToOrderingAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_operation_to_ordering_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWOperationToOrderingByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_operation_to_ordering_filter_str/{fs}/{fs2}",
		WrapAuth(GetWOperationToOrderingByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_operation_to_ordering_between_up_created_at/{fs}/{fs2}",
		WrapAuth(GetWOperationToOrderingBetweenUpCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_operation_to_product/{id:[0-9]+}",
		WrapAuth(GetWOperationToProduct, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_operation_to_product_get_all",
		WrapAuth(GetWOperationToProductAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_operation_to_product_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWOperationToProductByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_operation_to_product_filter_str/{fs}/{fs2}",
		WrapAuth(GetWOperationToProductByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_product_to_product/{id:[0-9]+}",
		WrapAuth(GetWProductToProduct, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_product_to_product_get_all",
		WrapAuth(GetWProductToProductAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_product_to_product_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWProductToProductByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_product_to_product_filter_str/{fs}/{fs2}",
		WrapAuth(GetWProductToProductByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_cbox_check/{id:[0-9]+}",
		WrapAuth(GetWCboxCheck, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_cbox_check_get_all",
		WrapAuth(GetWCboxCheckAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_cbox_check_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWCboxCheckByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_cbox_check_filter_str/{fs}/{fs2}",
		WrapAuth(GetWCboxCheckByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_cbox_check_between_created_at/{fs}/{fs2}",
		WrapAuth(GetWCboxCheckBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_item_to_cbox_check/{id:[0-9]+}",
		WrapAuth(GetWItemToCboxCheck, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_item_to_cbox_check_get_all",
		WrapAuth(GetWItemToCboxCheckAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_item_to_cbox_check_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWItemToCboxCheckByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_item_to_cbox_check_filter_str/{fs}/{fs2}",
		WrapAuth(GetWItemToCboxCheckByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_cash_in/{id:[0-9]+}",
		WrapAuth(GetWCashIn, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_cash_in_get_all",
		WrapAuth(GetWCashInAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_cash_in_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWCashInByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_cash_in_filter_str/{fs}/{fs2}",
		WrapAuth(GetWCashInByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_cash_in_between_created_at/{fs}/{fs2}",
		WrapAuth(GetWCashInBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_cash_out/{id:[0-9]+}",
		WrapAuth(GetWCashOut, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_cash_out_get_all",
		WrapAuth(GetWCashOutAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_cash_out_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWCashOutByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_cash_out_filter_str/{fs}/{fs2}",
		WrapAuth(GetWCashOutByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_cash_out_between_created_at/{fs}/{fs2}",
		WrapAuth(GetWCashOutBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_whs/{id:[0-9]+}",
		WrapAuth(GetWWhs, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_whs_get_all",
		WrapAuth(GetWWhsAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_whs_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWWhsByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_whs_filter_str/{fs}/{fs2}",
		WrapAuth(GetWWhsByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_whs_in/{id:[0-9]+}",
		WrapAuth(GetWWhsIn, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_whs_in_get_all",
		WrapAuth(GetWWhsInAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_whs_in_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWWhsInByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_whs_in_filter_str/{fs}/{fs2}",
		WrapAuth(GetWWhsInByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_whs_in_between_created_at/{fs}/{fs2}",
		WrapAuth(GetWWhsInBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_whs_in_between_contragent_created_at/{fs}/{fs2}",
		WrapAuth(GetWWhsInBetweenContragentCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_whs_out/{id:[0-9]+}",
		WrapAuth(GetWWhsOut, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_whs_out_get_all",
		WrapAuth(GetWWhsOutAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_whs_out_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWWhsOutByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_whs_out_filter_str/{fs}/{fs2}",
		WrapAuth(GetWWhsOutByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_whs_out_between_created_at/{fs}/{fs2}",
		WrapAuth(GetWWhsOutBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_whs_in/{id:[0-9]+}",
		WrapAuth(GetWMatherialToWhsIn, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_whs_in_get_all",
		WrapAuth(GetWMatherialToWhsInAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_whs_in_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWMatherialToWhsInByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_whs_in_filter_str/{fs}/{fs2}",
		WrapAuth(GetWMatherialToWhsInByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_whs_out/{id:[0-9]+}",
		WrapAuth(GetWMatherialToWhsOut, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_whs_out_get_all",
		WrapAuth(GetWMatherialToWhsOutAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_whs_out_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWMatherialToWhsOutByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_to_whs_out_filter_str/{fs}/{fs2}",
		WrapAuth(GetWMatherialToWhsOutByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_part/{id:[0-9]+}",
		WrapAuth(GetWMatherialPart, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_part_get_all",
		WrapAuth(GetWMatherialPartAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_part_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWMatherialPartByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_part_filter_str/{fs}/{fs2}",
		WrapAuth(GetWMatherialPartByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_part_between_created_at/{fs}/{fs2}",
		WrapAuth(GetWMatherialPartBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_part_slice/{id:[0-9]+}",
		WrapAuth(GetWMatherialPartSlice, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_part_slice_get_all",
		WrapAuth(GetWMatherialPartSliceAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_part_slice_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWMatherialPartSliceByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_part_slice_filter_str/{fs}/{fs2}",
		WrapAuth(GetWMatherialPartSliceByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_matherial_part_slice_between_created_at/{fs}/{fs2}",
		WrapAuth(GetWMatherialPartSliceBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_group/{id:[0-9]+}",
		WrapAuth(GetWProjectGroup, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_group_get_all",
		WrapAuth(GetWProjectGroupAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_group_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWProjectGroupByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_group_filter_str/{fs}/{fs2}",
		WrapAuth(GetWProjectGroupByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_status/{id:[0-9]+}",
		WrapAuth(GetWProjectStatus, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_status_get_all",
		WrapAuth(GetWProjectStatusAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_status_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWProjectStatusByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_status_filter_str/{fs}/{fs2}",
		WrapAuth(GetWProjectStatusByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_type/{id:[0-9]+}",
		WrapAuth(GetWProjectType, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_type_get_all",
		WrapAuth(GetWProjectTypeAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_type_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWProjectTypeByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_type_filter_str/{fs}/{fs2}",
		WrapAuth(GetWProjectTypeByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project/{id:[0-9]+}",
		WrapAuth(GetWProject, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_get_all",
		WrapAuth(GetWProjectAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWProjectByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_filter_str/{fs}/{fs2}",
		WrapAuth(GetWProjectByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_find_project_project_info_contragent_no_search_contact_no_search/{fs}",
		WrapAuth(GetWProjectFindByProjectInfoContragentNoSearchContactNoSearch, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_project_between_created_at/{fs}/{fs2}",
		WrapAuth(GetWProjectBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_counter/{id:[0-9]+}",
		WrapAuth(GetWCounter, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_counter_get_all",
		WrapAuth(GetWCounterAll, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_counter_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWCounterByFilterInt, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_counter_filter_str/{fs}/{fs2}",
		WrapAuth(GetWCounterByFilterStr, CATALOG_READ)).Methods("GET")

	r.HandleFunc("/w_record_to_counter/{id:[0-9]+}",
		WrapAuth(GetWRecordToCounter, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_record_to_counter_get_all",
		WrapAuth(GetWRecordToCounterAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_record_to_counter_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWRecordToCounterByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_record_to_counter_filter_str/{fs}/{fs2}",
		WrapAuth(GetWRecordToCounterByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_record_to_counter_between_created_at/{fs}/{fs2}",
		WrapAuth(GetWRecordToCounterBetweenCreatedAt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_wmc_number/{id:[0-9]+}",
		WrapAuth(GetWWmcNumber, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_wmc_number_get_all",
		WrapAuth(GetWWmcNumberAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_wmc_number_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWWmcNumberByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_wmc_number_filter_str/{fs}/{fs2}",
		WrapAuth(GetWWmcNumberByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_numbers_to_product/{id:[0-9]+}",
		WrapAuth(GetWNumbersToProduct, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_numbers_to_product_get_all",
		WrapAuth(GetWNumbersToProductAll, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_numbers_to_product_filter_int/{fs}/{id:[0-9]+}",
		WrapAuth(GetWNumbersToProductByFilterInt, DOC_READ)).Methods("GET")

	r.HandleFunc("/w_numbers_to_product_filter_str/{fs}/{fs2}",
		WrapAuth(GetWNumbersToProductByFilterStr, DOC_READ)).Methods("GET")

	r.HandleFunc("/upload/{id:[0-9]+}",
		WrapAuth(UploadFile, DOC_CREATE)).Methods("POST")

	r.HandleFunc("/product_to_ordering_default",
		WrapAuth(CreateProductToOrderingDefault, DOC_READ)).Methods("POST")

	r.HandleFunc("/product_deep/{id:[0-9]+}",
		WrapAuth(GetProductDeep, DOC_READ)).Methods("GET")

	r.HandleFunc("/product_complex/{id:[0-9]+}",
		WrapAuth(GetProductComplex, DOC_READ)).Methods("GET")

	r.HandleFunc("/project_dirs/{id:[0-9]+}",
		WrapAuth(CreateProjectDirs, DOC_CREATE)).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static/"))))

	r.HandleFunc("/login", WrapAuth(Login, LOGIN)).Methods("POST")
	r.HandleFunc("/logout", WrapAuth(Logout, LOGOUT)).Methods("GET")
	r.HandleFunc("/copy_project/{id:[0-9]+}", WrapAuth(CopyProject, ADMIN)).Methods("GET")
	r.HandleFunc("/copy_base/{fs}", WrapAuth(CopyBase, ADMIN)).Methods("GET")
	r.HandleFunc("/delete_base/{fs}", WrapAuth(DeleteBackupBase, ADMIN)).Methods("GET")
	r.HandleFunc("/get_bases", WrapAuth(GetBackupBases, ADMIN)).Methods("GET")
	r.HandleFunc("/restore_base/{fs}", WrapAuth(RestoreBaseFromBackup, ADMIN)).Methods("GET")
	r.HandleFunc("/ws", WrapAuth(UpgradeWS, WS_CONNECT)).Methods("GET")

	return r
}

type Config struct {
	Port          string   `json:"port"`
	BckpPath      string   `json:"bckp_path"`
	MaketsPath    string   `json:"makets_path"`
	OldMaketsPath string   `json:"old_makets_path"`
	NewMaketsPath string   `json:"new_makets_path"`
	MaketDirs     []string `json:"maket_dirs"`
	DBFile        string   `json:"db_file"`
}

var Cfg Config

func LoadConfig() error {
	data, err := os.ReadFile("config.json")
	if err == nil {
		decoder := json.NewDecoder(strings.NewReader(string(data)))
		err = decoder.Decode(&Cfg)
	}
	return err
}

func main() {

	if err := LoadConfig(); err != nil {
		log.Fatal("Can`t load config file!")
	}

	err := DBconnect(Cfg.DBFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("running on port ", Cfg.Port)
	log.Println("makets on path ", Cfg.MaketsPath)
	r := makeRouter()
	//log.Fatal(http.ListenAndServe("0.0.0.0:"+*portFlag, r))

	Root, err := tg.InitRoot("Target", tg.Horizontal)
	if err != nil {
		log.Fatal(err)
	}

	tg.SetTheme("clam")
	tg.AddStatusField()
	tg.AddStatus(0, "This is test start.")

	makeGui(Root, Cfg.Port, r)

	tg.MainLoop()
}

func makeGui(Root tg.Container, port string, router *mux.Router) {
	mainFrame := tg.NewSplitter(tg.Expand | tg.Horizontal)
	Root.Add(mainFrame)

	MainPanel := tg.NewBox(tg.Expand)
	RightPanel := tg.NewBox(0)
	mainFrame.Add(MainPanel, RightPanel)
	block := tg.NewButton("Start Server", 0)
	RightPanel.Add(block)
	block.IfPressed(func(s string) {
		tg.MessageBox("Starting server on port " + port)
		go http.ListenAndServe("0.0.0.0:"+port, router)
		block.Disable()
		tg.AddStatus(0, "Starting server on port "+port)
		go handleMessages()
	})
}
