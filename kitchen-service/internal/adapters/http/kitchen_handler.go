package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"kitchen-service/internal/ports"
)

type KitchenHandler struct {
	kitchenService ports.KitchenService
}

func NewKitchenHandler(kitchenService ports.KitchenService) *KitchenHandler {
	return &KitchenHandler{
		kitchenService: kitchenService,
	}
}

func (kh *KitchenHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/kitchen/orders", kh.GetAllOrders).Methods("GET")
	router.HandleFunc("/kitchen/orders/{orderID}", kh.GetOrder).Methods("GET")
	router.HandleFunc("/kitchen/orders/{orderID}/start", kh.StartPreparation).Methods("POST")
	router.HandleFunc("/kitchen/orders/{orderID}/complete", kh.CompleteOrder).Methods("POST")
	router.HandleFunc("/kitchen/orders/{orderID}/pickup", kh.MarkPickedUpByDriver).Methods("POST")
	router.HandleFunc("/kitchen/orders/{orderID}/cancel", kh.CancelOrder).Methods("POST")
	router.HandleFunc("/kitchen/stats", kh.GetKitchenStats).Methods("GET")
	router.HandleFunc("/kitchen/process-queue", kh.ProcessQueue).Methods("POST")
	router.HandleFunc("/kitchen/dashboard", kh.ServeDashboard).Methods("GET")
}

func (kh *KitchenHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := kh.kitchenService.GetAllOrders(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(orders)
}

func (kh *KitchenHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderID"]

	order, err := kh.kitchenService.GetOrderStatus(r.Context(), orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(order)
}

func (kh *KitchenHandler) StartPreparation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderID"]

	err := kh.kitchenService.StartPreparation(r.Context(), orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Preparation started",
		"orderID": orderID,
	})
}

func (kh *KitchenHandler) CompleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderID"]

	err := kh.kitchenService.CompleteOrder(r.Context(), orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Order completed",
		"orderID": orderID,
	})
}

func (kh *KitchenHandler) MarkPickedUpByDriver(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderID"]

	err := kh.kitchenService.MarkPickedUpByDriver(r.Context(), orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Order picked up by driver",
		"orderID": orderID,
	})
}

func (kh *KitchenHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderID"]

	err := kh.kitchenService.CancelOrder(r.Context(), orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Order cancelled",
		"orderID": orderID,
	})
}

func (kh *KitchenHandler) GetKitchenStats(w http.ResponseWriter, r *http.Request) {
	stats, err := kh.kitchenService.GetKitchenStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(stats)
}

func (kh *KitchenHandler) ProcessQueue(w http.ResponseWriter, r *http.Request) {
	err := kh.kitchenService.ProcessOrderQueue(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Queue processed successfully",
	})
}

// Kitchen dashboard will be integrated into React frontend
// TODO: Move kitchen management to React admin interface
func (kh *KitchenHandler) ServeDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"message": "Kitchen dashboard moved to React frontend"}`))
}