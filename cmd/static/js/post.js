
// Get modal elements
const signupModal = document.getElementById("signupModal");
const loginModal = document.getElementById("loginModal");

// Get buttons and close elements
const signupBtn = document.getElementById("signupBtn");
const loginBtn = document.getElementById("loginBtn");
const signupCloseBtn = document.getElementById("signupCloseBtn");
const loginCloseBtn = document.getElementById("loginCloseBtn");

const openLoginFromSignup = document.getElementById("openLoginFromSignup");
const openSignupFromLogin = document.getElementById("openSignupFromLogin");

// Show Signup Modal
signupBtn.onclick = function() {
    signupModal.style.display = "block";
}

// Show Login Modal
loginBtn.onclick = function() {
    loginModal.style.display = "block";
}

// Close Signup Modal
signupCloseBtn.onclick = function() {
    signupModal.style.display = "none";
}

// Close Login Modal
loginCloseBtn.onclick = function() {
    loginModal.style.display = "none";
}

// Close modals if clicked outside
window.onclick = function(event) {
    if (event.target == signupModal) {
        signupModal.style.display = "none";
    }
    if (event.target == loginModal) {
        loginModal.style.display = "none";
    }
}

// Switch to Login Modal from Signup Modal
openLoginFromSignup.onclick = function() {
    signupModal.style.display = "none";
    loginModal.style.display = "block";
}

// Switch to Signup Modal from Login Modal
openSignupFromLogin.onclick = function() {
    loginModal.style.display = "none";
    signupModal.style.display = "block";
}

// Simulate Signup Success
const signupForm = document.getElementById("signupForm");
signupForm.addEventListener("submit", function(event) {
    event.preventDefault(); // Simulate form submission

    alert("Signup successful! Redirecting to login page...");

    // Close signup modal and open login modal after signup
    signupModal.style.display = "none";
    loginModal.style.display = "block";
});

// Simulate Login Success
const loginForm = document.getElementById("loginForm");
loginForm.addEventListener("submit", function(event) {
    event.preventDefault(); // Simulate form submission

    alert("Login successful! Redirecting to homepage...");

    // Redirect to homepage after login
    window.location.href = "index.html"; // Replace with your actual homepage URL
});
