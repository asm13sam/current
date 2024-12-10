package main

type MatherialExtra struct {
	Matherial          WMatherial         `json:"matherial"`
	MatherialToProduct MatherialToProduct `json:"matherial_to_product"`
	Uid                int                `json:"uid"`
}

type OperationExtra struct {
	Operation          WOperation         `json:"operation"`
	OperationToProduct OperationToProduct `json:"operation_to_product"`
	Uid                int                `json:"uid"`
}

type ProductExtra struct {
	Product          WProduct         `json:"product"`
	ProductToProduct ProductToProduct `json:"product_to_product"`
	Uid              int              `json:"uid"`
}

type ProductComplex struct {
	Product         Product                     `json:"product"`
	MaterialToProd  map[string][]MatherialExtra `json:"matherial_extra"`
	OperationToProd map[string][]OperationExtra `json:"operation_extra"`
	ProductToProd   map[string][]ProductExtra   `json:"product_extra"`
	Uid             int                         `json:"uid"`
}

func ProductComplexGet(productId int) (ProductComplex, error) {
	var err error
	var pc ProductComplex
	counter := 1
	pc.Product, err = ProductGet(productId, nil)
	if err != nil {
		return pc, err
	}
	MaterialToProd, err := MatherialToProductGetByFilterInt("product_id", pc.Product.Id, false, false, nil)
	if err != nil {
		return pc, err
	}
	pc.MaterialToProd = make(map[string][]MatherialExtra)
	for _, m2p := range MaterialToProd {
		_, ok := pc.MaterialToProd[m2p.ListName]
		if !ok {
			pc.MaterialToProd[m2p.ListName] = []MatherialExtra{}
		}
		m, err := WMatherialGet(m2p.MatherialId)
		if err != nil {
			return pc, err
		}
		pc.MaterialToProd[m2p.ListName] = append(pc.MaterialToProd[m2p.ListName], MatherialExtra{m, m2p, counter})
		counter++
	}
	OperationToProd, err := OperationToProductGetByFilterInt("product_id", pc.Product.Id, false, false, nil)
	if err != nil {
		return pc, err
	}
	pc.OperationToProd = make(map[string][]OperationExtra)
	for _, o2p := range OperationToProd {
		_, ok := pc.OperationToProd[o2p.ListName]
		if !ok {
			pc.OperationToProd[o2p.ListName] = []OperationExtra{}
		}
		o, err := WOperationGet(o2p.OperationId)
		if err != nil {
			return pc, err
		}
		pc.OperationToProd[o2p.ListName] = append(pc.OperationToProd[o2p.ListName], OperationExtra{o, o2p, counter})
		counter++
	}
	ProductToProd, err := ProductToProductGetByFilterInt("product_id", pc.Product.Id, false, false, nil)
	if err != nil {
		return pc, err
	}
	pc.ProductToProd = make(map[string][]ProductExtra)
	for _, p2p := range ProductToProd {
		_, ok := pc.ProductToProd[p2p.ListName]
		if !ok {
			pc.ProductToProd[p2p.ListName] = []ProductExtra{}
		}
		p, err := WProductGet(p2p.Product2Id)
		if err != nil {
			return pc, err
		}
		pc.ProductToProd[p2p.ListName] = append(pc.ProductToProd[p2p.ListName], ProductExtra{p, p2p, counter})
		counter++
	}
	return pc, nil
}

type ProductDeep struct {
	ProductExtra    ProductExtra                `json:"product_extra"`
	MaterialToProd  map[string][]MatherialExtra `json:"matherial_extra"`
	OperationToProd map[string][]OperationExtra `json:"operation_extra"`
	ProductToProd   map[string][]ProductDeep    `json:"product_deep"`
	Uid             int                         `json:"uid"`
}

func ProductDeepGet(productId int, counter *int) (ProductDeep, error) {
	var err error
	var pd ProductDeep
	pd.ProductExtra.ProductToProduct = ProductToProduct{}
	pd.ProductExtra.Product, err = WProductGet(productId)
	if err != nil {
		return pd, err
	}
	pd.Uid = *counter
	*counter++
	pd.ProductExtra.Uid = *counter
	*counter++
	MaterialToProd, err := MatherialToProductGetByFilterInt("product_id", pd.ProductExtra.Product.Id, false, false, nil)
	if err != nil {
		return pd, err
	}
	pd.MaterialToProd = make(map[string][]MatherialExtra)
	for _, m2p := range MaterialToProd {
		_, ok := pd.MaterialToProd[m2p.ListName]
		if !ok {
			pd.MaterialToProd[m2p.ListName] = []MatherialExtra{}
		}
		m, err := WMatherialGet(m2p.MatherialId)
		if err != nil {
			return pd, err
		}
		pd.MaterialToProd[m2p.ListName] = append(pd.MaterialToProd[m2p.ListName], MatherialExtra{m, m2p, *counter})
		*counter++
	}
	OperationToProd, err := OperationToProductGetByFilterInt("product_id", pd.ProductExtra.Product.Id, false, false, nil)
	if err != nil {
		return pd, err
	}
	pd.OperationToProd = make(map[string][]OperationExtra)
	for _, o2p := range OperationToProd {
		_, ok := pd.OperationToProd[o2p.ListName]
		if !ok {
			pd.OperationToProd[o2p.ListName] = []OperationExtra{}
		}
		o, err := WOperationGet(o2p.OperationId)
		if err != nil {
			return pd, err
		}
		pd.OperationToProd[o2p.ListName] = append(pd.OperationToProd[o2p.ListName], OperationExtra{o, o2p, *counter})
		*counter++
	}
	ProductToProd, err := ProductToProductGetByFilterInt("product_id", pd.ProductExtra.Product.Id, false, false, nil)
	if err != nil {
		return pd, err
	}
	pd.ProductToProd = make(map[string][]ProductDeep)
	for _, p2p := range ProductToProd {
		_, ok := pd.ProductToProd[p2p.ListName]
		if !ok {
			pd.ProductToProd[p2p.ListName] = []ProductDeep{}
		}
		p, err := ProductGet(p2p.Product2Id, nil)
		if err != nil {
			return pd, err
		}
		child_pd, err := ProductDeepGet(p.Id, counter)
		if err != nil {
			return pd, err
		}
		child_pd.ProductExtra.ProductToProduct = p2p
		pd.ProductToProd[p2p.ListName] = append(pd.ProductToProd[p2p.ListName], child_pd)
	}
	return pd, nil
}

func ProductToOrderingCreateDefault(p ProductToOrdering) (ProductToOrdering, error) {
	p, err := ProductToOrderingCreate(p, nil)
	if err != nil {
		return p, err
	}
	m2ps, err := MatherialToProductGetByFilterInt("product_id", p.ProductId, false, false, nil)
	if err != nil {
		return p, err
	}
	for _, m2p := range m2ps {
		if m2p.ListName != "default" {
			continue
		}
		m2o := MatherialToOrdering{
			Id:                  0,
			OrderingId:          p.OrderingId,
			MatherialId:         m2p.MatherialId,
			Width:               0.0,
			Length:              0.0,
			Pieces:              1,
			ColorId:             0,
			UserId:              p.UserId,
			Number:              p.Number * m2p.Number,
			Price:               m2p.Cost,
			Persent:             0.0,
			Profit:              0.0,
			Cost:                p.Number * m2p.Cost,
			Comm:                "",
			ProductToOrderingId: p.Id,
			IsActive:            true,
		}
		_, err = MatherialToOrderingCreate(m2o, nil)
		if err != nil {
			return p, err
		}
	}
	o2ps, err := OperationToProductGetByFilterInt("product_id", p.ProductId, false, false, nil)
	if err != nil {
		return p, err
	}
	for _, o2p := range o2ps {
		if o2p.ListName != "default" {
			continue
		}
		o, err := OperationGet(o2p.OperationId, nil)
		if err != nil {
			return p, err
		}
		o2o := OperationToOrdering{
			Id:                  0,
			OrderingId:          p.OrderingId,
			OperationId:         o2p.OperationId,
			UserId:              p.UserId,
			Number:              p.Number * o2p.Number,
			Price:               o2p.Cost,
			UserSum:             o.Price * p.Number,
			Cost:                p.Number * o2p.Cost,
			EquipmentId:         o2p.EquipmentId,
			EquipmentCost:       o2p.EquipmentCost * p.Number,
			Comm:                "",
			ProductToOrderingId: p.Id,
			IsActive:            true,
		}
		_, err = OperationToOrderingCreate(o2o, nil)
		if err != nil {
			return p, err
		}
	}
	return p, nil
}
