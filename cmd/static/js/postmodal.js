// Modal handling
const createPostModal = document.getElementById('createPostModal');
const createPostBtn = document.getElementById('createPostBtn');
const createPostCloseBtn = document.getElementById('createPostCloseBtn');
const createPostForm = document.getElementById('createPostForm');

// Show create post modal
createPostBtn.addEventListener('click', () => {
    createPostModal.style.display = 'block';
});

// Close create post modal
createPostCloseBtn.addEventListener('click', () => {
    createPostModal.style.display = 'none';
});

// Close modal when clicking outside
window.addEventListener('click', (e) => {
    if (e.target === createPostModal) {
        createPostModal.style.display = 'none';
    }
});

// Handle post creation
createPostForm.addEventListener('submit', async (e) => {
    // e.preventDefault();

    const title = createPostForm.querySelector('input[type="text"]').value;
    const content = createPostForm.querySelector('textarea').value;
    const categoryInputs = createPostForm.querySelectorAll('input[name="categories"]:checked');
    
    // Convert category IDs to integers
    const categories = Array.from(categoryInputs).map(input => parseInt(input.value, 10));

    // Validate form
    if (!title.trim()) {
        alert('Please enter a title');
        return;
    }
    if (!content.trim()) {
        alert('Please enter content');
        return;
    }
    if (categories.length === 0) {
        alert('Please select at least one category');
        return;
    }

    try {
        const response = await fetch('/api/posts/create', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                title: title.trim(),
                content: content.trim(),
                categories: categories
            })
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || 'Failed to create post');
        }
    
        const post = await response.json();
    
        // Clear form and close modal
        createPostForm.reset();
        createPostModal.style.display = 'none';
    
        // Refresh the page to show the new post
        window.location.reload();
    } catch (error) {
        console.error('Error creating post:', error);
        alert(error.message || 'Failed to create post. Please try again.');
    }
});

// Load categories when the page loads
async function loadCategories() {
    try {
        const response = await fetch('/api/categories');
        if (!response.ok) {
            throw new Error('Failed to fetch categories');
        }
        const categories = await response.json();
        
        // Update category checkboxes in create post form
        const categoryContainer = document.querySelector('.category-selection');
        if (!categoryContainer) return;
        
        categoryContainer.innerHTML = '<h3>Select Categories</h3>';
        
        categories.forEach(category => {
            const label = document.createElement('label');
            label.className = 'category-checkbox';

            // Display the selected categories
            label.innerHTML = `
                <input type="checkbox" name="categories" value="${category.id}">
                <span>${category.name}</span>
            `;
            categoryContainer.appendChild(label);
        });
    } catch (error) {
        console.error('Error loading categories:', error);
    }
}

// Initialize categories when page loads
document.addEventListener('DOMContentLoaded', loadCategories);
