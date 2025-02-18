const originalFetch = window.fetch;
window.fetch = function(url, options = {}) {
    const token = localStorage.getItem('token');

    if (token) {
        options.headers = options.headers || {};
        if (options.headers instanceof Headers) {
            options.headers.append('Authorization', `Bearer ${token}`);
        } else {
            options.headers['Authorization'] = `Bearer ${token}`;
        }
    }

    return originalFetch(url, options);
};

document.addEventListener("DOMContentLoaded", () => {
    const loginForm = document.getElementById("loginForm");

    if (loginForm) {
        loginForm.addEventListener("submit", async (e) => {
            e.preventDefault();

            const username = document.getElementById("username").value;
            const password = document.getElementById("password").value;

            const response = await fetch("http://localhost:8080/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({ username, password })
            });

            const data = await response.json();

            if (response.ok) {
                localStorage.setItem("token", data.token);
                alert("Login successful!");
                window.location.href = "/";
            } else {
                alert(data.error || "Login failed");
            }
        });
    }

    const registerForm = document.getElementById("registerForm");

    if (registerForm) {
        registerForm.addEventListener("submit", async (e) => {
            e.preventDefault();

            const username = document.getElementById("username").value;
            const password = document.getElementById("password").value;
            const confirmPassword = document.getElementById("confirmPassword").value;

            if (password !== confirmPassword) {
                alert("Passwords do not match!");
                return;
            }

            try {
                const response = await fetch("http://localhost:8080/register", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({ username, password })
                });

                const data = await response.json();

                if (response.ok) {
                    alert("Registration successful! Please log in.");
                    window.location.href = "/front/login";
                } else {
                    alert(data.error || "Registration failed");
                }
            } catch (error) {
                console.error("Error during registration:", error);
                alert("An error occurred during registration");
            }
        });
    }
});

document.addEventListener("DOMContentLoaded", () => {
    const registerForm = document.getElementById("registerForm");

    if (registerForm) {
        registerForm.addEventListener("submit", async (e) => {
            e.preventDefault();

            const username = document.getElementById("username").value;
            const password = document.getElementById("password").value;
            const confirmPassword = document.getElementById("confirmPassword").value;

            if (password !== confirmPassword) {
                alert("Passwords do not match!");
                return;
            }

            console.log("Attempting to register:", { username, password });

            try {
                const response = await fetch("http://localhost:8080/register", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({ username, password })
                });

                console.log("Response status:", response.status);

                const data = await response.json();
                console.log("Response data:", data);

                if (response.ok) {
                    alert("Registration successful! Please log in.");
                    window.location.href = "/front/login";
                } else {
                    alert(data.error || "Registration failed");
                }
            } catch (error) {
                console.error("Error during registration:", error);
                alert("An error occurred during registration");
            }
        });
    }
});
