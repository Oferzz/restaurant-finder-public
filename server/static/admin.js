// Login Form Submission Handler
document.getElementById('login-form').addEventListener('submit', async (e) => {
    e.preventDefault(); // Prevent default form submission behavior

    const password = document.getElementById('password').value;

    try {
        // Perform a dummy request to validate the admin password
        const response = await fetch('/admin/audit-logs', {
            method: 'GET',
            headers: { 'Authorization': password },
        });

        if (response.ok) {
            // Hide the login section and display the admin section
            document.getElementById('login-section').style.display = 'none';
            document.getElementById('admin-section').style.display = 'block';

            // Save the password in localStorage for later use
            localStorage.setItem('admin-password', password);
        } else {
            alert('Invalid password');
        }
    } catch (error) {
        console.error('Error during login:', error);
        alert('Failed to log in.');
    }
});

// Add Restaurant Form Submission Handler
document.getElementById('add-restaurant-form').addEventListener('submit', async (e) => {
    e.preventDefault(); // Prevent default form submission behavior

    const restaurant = {
        restaurant_name: document.getElementById('restaurant_name').value,
        address: document.getElementById('address').value,
        phone: document.getElementById('phone').value,
        website: document.getElementById('website').value,
        cuisine_type: document.getElementById('cuisine_type').value,
        is_kosher: document.getElementById('is_kosher').checked,
        opening_hours: {
            Monday: document.getElementById('monday').value,
            Tuesday: document.getElementById('tuesday').value,
            Wednesday: document.getElementById('wednesday').value,
            Thursday: document.getElementById('thursday').value,
            Friday: document.getElementById('friday').value,
            Saturday: document.getElementById('saturday').value,
            Sunday: document.getElementById('sunday').value,
        },
    };

    const password = localStorage.getItem('admin-password'); // Retrieve stored password

    try {
        // Send POST request to add a new restaurant
        const response = await fetch('/admin/restaurants', {
            method: 'POST',
            headers: {
                'Authorization': password,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(restaurant),
        });

        if (response.ok) {
            alert('Restaurant added successfully');
            document.getElementById('add-restaurant-form').reset(); // Clear form after successful submission
        } else {
            alert('Failed to add restaurant');
        }
    } catch (error) {
        console.error('Error adding restaurant:', error);
        alert('Error occurred while adding the restaurant.');
    }
});

// Fetch Restaurant Details by ID
document.getElementById('fetch-restaurant-btn').addEventListener('click', async () => {
    const restaurantID = document.getElementById('edit_restaurant_id').value;
    const password = localStorage.getItem('admin-password'); // Retrieve admin password

    if (!restaurantID) {
        alert('Please enter a Restaurant ID.');
        return;
    }

    try {
        const response = await fetch(`/admin/restaurants/${restaurantID}`, {
            headers: { 'Authorization': password },
        });

        if (response.ok) {
            const restaurant = await response.json();

            // Populate the form fields with the retrieved data
            document.getElementById('edit_restaurant_name').value = restaurant.restaurant_name || '';
            document.getElementById('edit_address').value = restaurant.address || '';
            document.getElementById('edit_phone').value = restaurant.phone || '';
            document.getElementById('edit_website').value = restaurant.website || '';
            document.getElementById('edit_cuisine_type').value = restaurant.cuisine_type || '';
            document.getElementById('edit_is_kosher').checked = restaurant.is_kosher || false;

            const openingHours = restaurant.opening_hours || {};
            document.getElementById('edit_monday').value = openingHours.Monday || '';
            document.getElementById('edit_tuesday').value = openingHours.Tuesday || '';
            document.getElementById('edit_wednesday').value = openingHours.Wednesday || '';
            document.getElementById('edit_thursday').value = openingHours.Thursday || '';
            document.getElementById('edit_friday').value = openingHours.Friday || '';
            document.getElementById('edit_saturday').value = openingHours.Saturday || '';
            document.getElementById('edit_sunday').value = openingHours.Sunday || '';

            // Show the edit form
            document.getElementById('edit-restaurant-form').style.display = 'block';
        } else {
            alert('Failed to fetch restaurant details. Ensure the Restaurant ID is correct.');
        }
    } catch (error) {
        console.error('Error fetching restaurant details:', error);
        alert('Error occurred while fetching the restaurant details.');
    }
});

// Edit Restaurant Form Submission Handler
document.getElementById('edit-restaurant-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const restaurantID = document.getElementById('edit_restaurant_id').value;
    const restaurant = {
        restaurant_name: document.getElementById('edit_restaurant_name').value,
        address: document.getElementById('edit_address').value,
        phone: document.getElementById('edit_phone').value,
        website: document.getElementById('edit_website').value,
        cuisine_type: document.getElementById('edit_cuisine_type').value,
        is_kosher: document.getElementById('edit_is_kosher').checked,
        opening_hours: {
            Monday: document.getElementById('edit_monday').value,
            Tuesday: document.getElementById('edit_tuesday').value,
            Wednesday: document.getElementById('edit_wednesday').value,
            Thursday: document.getElementById('edit_thursday').value,
            Friday: document.getElementById('edit_friday').value,
            Saturday: document.getElementById('edit_saturday').value,
            Sunday: document.getElementById('edit_sunday').value,
        },
    };

    const password = localStorage.getItem('admin-password'); // Retrieve admin password

    try {
        const response = await fetch(`/admin/restaurants/${restaurantID}`, {
            method: 'PUT',
            headers: {
                'Authorization': password,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(restaurant),
        });

        if (response.ok) {
            alert('Restaurant updated successfully');
            document.getElementById('edit-restaurant-form').reset(); // Clear form after submission
            document.getElementById('edit-restaurant-form').style.display = 'none'; // Hide form
        } else {
            alert('Failed to update restaurant');
        }
    } catch (error) {
        console.error('Error updating restaurant:', error);
        alert('Error occurred while updating the restaurant.');
    }
});

// Remove Restaurant Form Submission Handler
document.getElementById('remove-restaurant-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const restaurantID = document.getElementById('remove_restaurant_id').value;
    const password = localStorage.getItem('admin-password');

    try {
        const response = await fetch(`/admin/restaurants/${restaurantID}`, {
            method: 'DELETE',
            headers: { 'Authorization': password },
        });

        if (response.ok) {
            alert('Restaurant removed successfully');
            document.getElementById('remove-restaurant-form').reset();
        } else {
            alert('Failed to remove restaurant');
        }
    } catch (error) {
        console.error('Error removing restaurant:', error);
        alert('Error occurred while removing the restaurant.');
    }
});

// Fetch Audit Logs Button Handler
document.getElementById('fetch-audit-logs-btn').addEventListener('click', async () => {
    console.log("Audit logs button clicked"); // Debugging
    const password = localStorage.getItem('admin-password'); // Retrieve stored password

    if (!password) {
        alert('Please log in first.');
        return;
    }

    try {
        // Send GET request to fetch audit logs
        const response = await fetch('/admin/audit-logs', {
            headers: { 'Authorization': password },
        });

        if (response.ok) {
            const logs = await response.json();
            const tbody = document.getElementById('audit-log-table').querySelector('tbody');
            tbody.innerHTML = ''; // Clear previous logs

            // Populate the audit log table with new data
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

            // Show the table
            document.getElementById('audit-log-table').style.display = 'table';
        } else {
            alert('Failed to fetch audit logs.');
        }
    } catch (error) {
        console.error('Error fetching audit logs:', error);
        alert('Failed to fetch audit logs.');
    }
});