// Simple JavaScript to handle form submissions and call API endpoints

document.getElementById('loginForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const data = {
        account: document.getElementById('loginAccount').value,
        password: document.getElementById('loginPassword').value,
        school_id: document.getElementById('loginSchool').value
    };
    const resp = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });
    const res = await resp.json();
    document.getElementById('loginResult').textContent = res.message || res.error;
});

document.getElementById('queryForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const account = document.getElementById('queryAccount').value;
    const resp = await fetch('/api/query?account=' + encodeURIComponent(account));
    const res = await resp.json();
    document.getElementById('queryResult').textContent = JSON.stringify(res, null, 2);
});

document.getElementById('signForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const data = {
        account: document.getElementById('signAccount').value,
        address: document.getElementById('signAddress').value,
        address_name: document.getElementById('signAddressName').value,
        latitude: document.getElementById('signLatitude').value,
        longitude: document.getElementById('signLongitude').value
    };
    const resp = await fetch('/api/sign', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });
    const res = await resp.json();
    document.getElementById('signResult').textContent = res.message || res.error;
});
