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

func (kh *KitchenHandler) ServeDashboard(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="de">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>🍳 Kitchen Dashboard</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: 'Arial', sans-serif; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); min-height: 100vh; }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; color: white; margin-bottom: 30px; }
        .header h1 { font-size: 2.5em; margin-bottom: 10px; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .stat-card { background: white; padding: 20px; border-radius: 10px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; }
        .stat-number { font-size: 2em; font-weight: bold; color: #667eea; }
        .stat-label { color: #666; margin-top: 5px; }
        .orders-section { background: white; padding: 20px; border-radius: 10px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        .orders-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
        .btn { background: #667eea; color: white; border: none; padding: 10px 20px; border-radius: 5px; cursor: pointer; font-size: 14px; }
        .btn:hover { background: #5a67d8; }
        .orders-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 20px; }
        .order-card { border: 1px solid #e2e8f0; border-radius: 8px; padding: 15px; background: #f8f9fa; }
        .order-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
        .order-id { font-weight: bold; color: #2d3748; }
        .status { padding: 4px 12px; border-radius: 20px; font-size: 12px; font-weight: bold; }
        .status.received { background: #fef5e7; color: #744210; }
        .status.preparing { background: #e6fffa; color: #065f46; }
        .status.ready { background: #dcfce7; color: #166534; }
        .status.picked_up_by_driver { background: #ddd6fe; color: #5b21b6; }
        .status.cancelled { background: #fee2e2; color: #991b1b; }
        .order-items { margin: 10px 0; font-size: 14px; color: #4a5568; }
        .order-actions { display: flex; gap: 5px; flex-wrap: wrap; margin-top: 10px; }
        .btn-small { padding: 5px 10px; font-size: 12px; border: none; border-radius: 4px; cursor: pointer; }
        .btn-start { background: #48bb78; color: white; }
        .btn-complete { background: #4299e1; color: white; }
        .btn-pickup { background: #9f7aea; color: white; }
        .btn-cancel { background: #f56565; color: white; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🍳 Kitchen Dashboard</h1>
            <p>Manage your restaurant orders in real-time</p>
        </div>
        
        <div class="stats" id="stats">
            <!-- Stats will be loaded here -->
        </div>
        
        <div class="orders-section">
            <div class="orders-header">
                <h2>📋 Active Orders</h2>
                <button class="btn" onclick="processQueue()">🚀 Process Queue</button>
                <button class="btn" onclick="refreshData()">🔄 Refresh</button>
            </div>
            <div class="orders-grid" id="orders">
                <!-- Orders will be loaded here -->
            </div>
        </div>
    </div>

    <script>
        let ordersData = [];

        async function loadStats() {
            try {
                const response = await fetch('/kitchen/stats');
                const stats = await response.json();
                
                const statsHtml = ` + "`" + `
                    <div class="stat-card">
                        <div class="stat-number">${stats.total_orders || 0}</div>
                        <div class="stat-label">Total Orders</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-number">${stats.orders_received || 0}</div>
                        <div class="stat-label">📥 Received</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-number">${stats.orders_preparing || 0}</div>
                        <div class="stat-label">👨‍🍳 Preparing</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-number">${stats.orders_ready || 0}</div>
                        <div class="stat-label">✅ Ready</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-number">${stats.orders_picked_up_by_driver || 0}</div>
                        <div class="stat-label">🚗 Picked Up</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-number">${stats.average_wait_time || 0} sec</div>
                        <div class="stat-label">⏱️ Avg Wait (sec)</div>
                    </div>
                ` + "`" + `;
                
                document.getElementById('stats').innerHTML = statsHtml;
            } catch (error) {
                console.error('Error loading stats:', error);
            }
        }

        async function loadOrders() {
            try {
                const response = await fetch('/kitchen/orders');
                ordersData = await response.json() || [];
                renderOrders();
            } catch (error) {
                console.error('Error loading orders:', error);
                document.getElementById('orders').innerHTML = '<p>Error loading orders</p>';
            }
        }

        function renderOrders() {
            const ordersHtml = ordersData.map(order => ` + "`" + `
                <div class="order-card">
                    <div class="order-header">
                        <span class="order-id">Order #${order.order_id}</span>
                        <span class="status ${order.status}">${order.status}</span>
                    </div>
                    <div class="order-items">
                        <strong>Items:</strong> ${order.items?.map(item => ` + "`" + `${item.quantity}x ${item.product_name}` + "`" + `).join(', ') || 'N/A'}
                    </div>
                    <div class="order-items">
                        <strong>Customer:</strong> ${order.customer_id}<br>
                        <strong>Est. Time:</strong> ${order.estimated_time} sec
                    </div>
                    ${getOrderActions(order)}
                </div>
            ` + "`" + `).join('');
            
            document.getElementById('orders').innerHTML = ordersHtml || '<p>No orders found</p>';
        }

        function getOrderActions(order) {
            let actions = '';
            
            switch(order.status) {
                case 'received':
                    actions = ` + "`" + `<button class="btn-small btn-start" onclick="startPreparation('${order.order_id}')">🔥 Start Prep</button>
                               <button class="btn-small btn-cancel" onclick="cancelOrder('${order.order_id}')">❌ Cancel</button>` + "`" + `;
                    break;
                case 'preparing':
                    actions = ` + "`" + `<button class="btn-small btn-complete" onclick="completeOrder('${order.order_id}')">✅ Complete</button>
                               <button class="btn-small btn-cancel" onclick="cancelOrder('${order.order_id}')">❌ Cancel</button>` + "`" + `;
                    break;
                case 'ready':
                    actions = ` + "`" + `<button class="btn-small btn-pickup" onclick="markPickedUpByDriver('${order.order_id}')">🚗 Driver Pickup</button>` + "`" + `;
                    break;
            }
            
            return actions ? ` + "`" + `<div class="order-actions">${actions}</div>` + "`" + ` : '';
        }

        async function startPreparation(orderID) {
            try {
                const response = await fetch(` + "`" + `/kitchen/orders/${orderID}/start` + "`" + `, { method: 'POST' });
                if (response.ok) {
                    refreshData();
                } else {
                    alert('Error starting preparation');
                }
            } catch (error) {
                console.error('Error:', error);
            }
        }

        async function completeOrder(orderID) {
            try {
                const response = await fetch(` + "`" + `/kitchen/orders/${orderID}/complete` + "`" + `, { method: 'POST' });
                if (response.ok) {
                    refreshData();
                } else {
                    alert('Error completing order');
                }
            } catch (error) {
                console.error('Error:', error);
            }
        }

        async function markPickedUpByDriver(orderID) {
            try {
                const response = await fetch(` + "`" + `/kitchen/orders/${orderID}/pickup` + "`" + `, { method: 'POST' });
                if (response.ok) {
                    refreshData();
                } else {
                    alert('Error marking order as picked up');
                }
            } catch (error) {
                console.error('Error:', error);
            }
        }

        async function cancelOrder(orderID) {
            if (confirm('Are you sure you want to cancel this order?')) {
                try {
                    const response = await fetch(` + "`" + `/kitchen/orders/${orderID}/cancel` + "`" + `, { method: 'POST' });
                    if (response.ok) {
                        refreshData();
                    } else {
                        alert('Error cancelling order');
                    }
                } catch (error) {
                    console.error('Error:', error);
                }
            }
        }

        async function processQueue() {
            try {
                const response = await fetch('/kitchen/process-queue', { method: 'POST' });
                if (response.ok) {
                    alert('Queue processed successfully!');
                    refreshData();
                } else {
                    alert('Error processing queue');
                }
            } catch (error) {
                console.error('Error:', error);
            }
        }

        function refreshData() {
            loadStats();
            loadOrders();
        }

        // Initial load and auto-refresh
        refreshData();
        setInterval(refreshData, 10000); // Refresh every 10 seconds
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}