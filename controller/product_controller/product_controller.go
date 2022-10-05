package product_controller

import (
	"go_jwt_mux/helper"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	data := []map[string]interface{}{
		{
			"id":          1,
			"nama_produk": "kemeja",
			"stok":        100,
		},
		{
			"id":          2,
			"nama_produk": "celana",
			"stok":        200,
		},
		{
			"id":          1,
			"nama_produk": "sepatu",
			"stok":        220,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, data)
}
