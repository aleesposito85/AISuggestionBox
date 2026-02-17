// script.js

// DOM Elements
const form = document.getElementById('suggestionForm');
const suggestionList = document.getElementById('suggestionList');
const emptyState = document.getElementById('emptyState');
const countDisplay = document.getElementById('countDisplay');

// State Management
let suggestions = [];

// Initialize App
document.addEventListener('DOMContentLoaded', () => {
    loadSuggestions();
});

// Event Listenersl
form.addEventListener('submit', handleFormSubmit);

// Functions
async function handleFormSubmit(e) {
    e.preventDefault();

    // Get values
    const category = document.getElementById('category').value;
    const name = document.getElementById('name').value.trim();
    const message = document.getElementById('message').value.trim();
    const email = document.getElementById('email').value.trim();

    if (!email || !message) return;

    const newSuggestion = {
        category,
        name,
        email,
        message
    };

    try {
        const response = await fetch('https://suggestionboxapi.aleespohome.com/api/suggestions', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(newSuggestion)
        });

        if (response.ok) {
            const data = await response.json();
            console.log("Server response:", data);
            
            // Add the new suggestion locally for immediate UI feedback
            newSuggestion.id = data.id || Date.now(); // Use ID from DB if available
            newSuggestion.date = data.date
            suggestions.unshift(newSuggestion);
            
            renderSuggestions();
            form.reset();
            showToast();
        } else {
            alert('Error submitting suggestion. Please try again.');
        }
    } catch (error) {
        console.error('Fetch error:', error);
        alert('Connection error. Make sure the Go server is running.');
    }
}

async function loadSuggestions() {
    try {
        const response = await fetch('https://suggestionboxapi.aleespohome.com/api/suggestions');
        
        if (response.ok) {
            suggestions = await response.json();
            renderSuggestions();
        } else {
            console.error("Failed to load suggestions");
        }
    } catch (error) {
        console.error('Fetch error:', error);
        // Optional: Show error in UI if DB is down
    }
}

function toggleReply(id) {
    const replySection = document.getElementById(`reply-${id}`);
    const replyContent = document.getElementById(`reply-content-${id}`);
    const btn = document.getElementById(`reply-btn-${id}`);
    
    // Toggle class 'active' to show/hide
    replySection.classList.toggle('active');
    
    // Update button text
    if (replySection.classList.contains('active')) {
        btn.textContent = "Hide";
        replyContent.style.display = "block";
    } else {
        btn.textContent = "View Reply";
        replyContent.style.display = "none";
    }
}

function renderSuggestions() {
    // Toggle Empty State
    if (suggestions.length === 0) {
        emptyState.style.display = 'block';
        suggestionList.style.display = 'none';
        countDisplay.textContent = '0 Items';
        return;
    }

    emptyState.style.display = 'none';
    suggestionList.style.display = 'grid';
    countDisplay.textContent = `${suggestions.length} Item${suggestions.length !== 1 ? 's' : ''}`;

    // Clear list
    suggestionList.innerHTML = '';

    // Loop and Create HTML
    suggestions.forEach(suggestion => {
        const card = document.createElement('article');
        card.className = 'suggestion-card';
        
        // Determine badge class
        let badgeClass = 'badge-general';
        let categoryLabel = 'General';
        if (suggestion.category === 'feature') {
            badgeClass = 'badge-feature';
            categoryLabel = 'Feature';
        } else if (suggestion.category === 'bug') {
            badgeClass = 'badge-bug';
            categoryLabel = 'Bug';
        }

        card.innerHTML = `
            <div class="card-header">
                <span class="category-badge ${badgeClass}">${categoryLabel}</span>
            </div>
            <h3 class="card-title">${escapeHtml(suggestion.name)}</h3>
            <p class="card-desc">${escapeHtml(suggestion.message)}</p>
            <div class="reply-section" id="reply-${suggestion.id}">
                <div class="reply-header">
                    <button id="reply-btn-${suggestion.id}" class="badge-general category-badge reply-btn" onclick="toggleReply(${suggestion.id})">View Reply</button>
                </div>
                <div id="reply-content-${suggestion.id}" class="reply-content" style="display: none">
                    <div style="display:flex; align-items:center; color: var(--primary); font-weight:600;">
                        Team (AI) Reply
                    </div>
                    ${escapeHtml(suggestion.aiReply || "No reply provided yet.")}
                </div>
            </div>
            <div class="card-footer">
                <span>Submitted: ${new Date(suggestion.date).toLocaleDateString()}</span>
            </div>
        `;
        suggestionList.appendChild(card);
    });
}

function showToast() {
    const toast = document.getElementById('toast');
    toast.classList.add('show');
    setTimeout(() => {
        toast.classList.remove('show');
    }, 3000);
}

// Security Helper: Prevent XSS
function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return text.replace(/[&<>"']/g, function(m) { return map[m]; });
}