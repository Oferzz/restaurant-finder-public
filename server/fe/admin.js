document.getElementById('login-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const password = document.getElementById('password').value;

    try {
        const response = await fetch('/admin/audit-logs', {
            method: 'GET',
            headers: { 'Authorization': password },
        });

        if (response.ok) {
            document.getElementById('login-section').style.display = 'none';
            document.getElementById('admin-section').style.display = 'block';
            localStorage.setItem('admin-password', password);
            loadRestaurants(password);
            loadAuditLogs(password);
        } else {
            alert('Invalid password');
        }
    } catch (error) {
        console.error('Error during login:', error);
        alert('Failed to log in.');
    }
});

async function loadRestaurants(password) {
    try {
        const response = await fetch('/restaurants', {
            headers: { 'Authorization': password },
        });
        const restaurants = await response.json();
        const tbody = document.getElementById('restaurant-table').querySelector('tbody');
        tbody.innerHTML = '';
        restaurants.forEach((restaurant) => {
            const row = `
                <tr>
                    <td>${restaurant.restaurant_id}</td>
                    <td>${restaurant.restaurant_name}</td>
                    <td>
                        <button onclick="editRestaurant('${restaurant.restaurant_id}', '${password}')">Edit</button>
                        <button onclick="deleteRestaurant('${restaurant.restaurant_id}', '${password}')">Delete</button>
                    </td>
                </tr>
            `;
            tbody.innerHTML += row;
        });
    } catch (error) {
        console.error('Error loading restaurants:', error);
        alert('Failed to load restaurants.');
    }
}

async function loadAuditLogs(password) {
    try {
        const response = await fetch('/admin/audit-logs', {
            headers: { 'Authorization': password },
        });
        const logs = await response.json();
        const tbody = document.getElementById('audit-log-table').querySelector('tbody');
        tbody.innerHTML = '';
        logs.forEach((log) => {
            const row = `
                <tr>
                    <td>${log.timestamp}</td>
                    <td>${log.query}</td>
                    <td>${log.ip}</td>
                    <td>${log.country}</td>
                </tr>
            `;
            tbody.innerHTML += row;
        });
    } catch (error) {
        console.error('Error loading audit logs:', error);
        alert('Failed to load audit logs.');
    }
}

async function deleteRestaurant(id, password) {
    try {
        const response = await fetch(`/admin/restaurants/${id}`, {
            method: 'DELETE',
            headers: { 'Authorization': password },
        });
        if (response.ok) {
            alert('Restaurant deleted successfully.');
            loadRestaurants(password);
        } else {
            alert('Failed to delete restaurant.');
        }
    } catch (error) {
        console.error('Error deleting restaurant:', error);
        alert('Failed to delete restaurant.');
    }
}

function editRestaurant(id, password) {
    alert(`Editing restaurant: ${id}`); // Implement edit functionality here
}
