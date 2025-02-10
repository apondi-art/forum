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
    e.preventDefault();

    const title = createPostForm.querySelector('input[type="text"]').value;
    const content = createPostForm.querySelector('textarea').value;
    const categoryInputs = createPostForm.querySelectorAll('input[name="categories"]:checked');
    const categories = Array.from(categoryInputs).map(input => parseInt(input.value));

    try {
        const response = await fetch('/api/posts/create', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                title,
                content,
                categories
            })
        });

        if (!response.ok) {
            throw new Error('Failed to create post');
        }

        // Refresh the page to show the new post
        window.location.reload();
    } catch (error) {
        console.error('Error creating post:', error);
        alert('Failed to create post. Please try again.');
    }
});

// // Like/Dislike handling
// document.addEventListener('click', async (e) => {
//     if (e.target.classList.contains('like-btn') || e.target.classList.contains('dislike-btn')) {
//         const button = e.target;
//         const postId = button.dataset.postId;
//         const type = button.classList.contains('like-btn') ? 'like' : 'dislike';

//         try {
//             const response = await fetch('/api/reaction', {
//                 method: 'POST',
//                 headers: {
//                     'Content-Type': 'application/json',
//                 },
//                 body: JSON.stringify({
//                     postId,
//                     type
//                 })
//             });

//             if (!response.ok) {
//                 throw new Error('Failed to update reaction');
//             }

//             const data = await response.json();
            
//             // Update like/dislike counts
//             const postCard = button.closest('.post-card');
//             postCard.querySelector('.like-count').textContent = data.likes;
//             postCard.querySelector('.dislike-count').textContent = data.dislikes;
//         } catch (error) {
//             console.error('Error updating reaction:', error);
//             alert('Failed to update reaction. Please try again.');
//         }
//     }
// });

// Comment section toggle
document.addEventListener('click', (e) => {
    if (e.target.classList.contains('view-comments')) {
        const postId = e.target.dataset.postId;
        const commentsSection = document.getElementById(`comments-${postId}`);
        commentsSection.style.display = commentsSection.style.display === 'none' ? 'block' : 'none';
    }
});

// // Comment form handling
// document.addEventListener('submit', async (e) => {
//     if (e.target.classList.contains('comment-form')) {
//         e.preventDefault();
//         const form = e.target;
//         const postId = form.dataset.postId;
//         const content = form.querySelector('textarea').value;

//         try {
//             const response = await fetch('/api/comment', {
//                 method: 'POST',
//                 headers: {
//                     'Content-Type': 'application/json',
//                 },
//                 body: JSON.stringify({
//                     postId,
//                     content
//                 })
//             });

//             if (!response.ok) {
//                 throw new Error('Failed to create comment');
//             }

//             const comment = await response.json();
            
//             // Add new comment to the page
//             const commentsSection = document.getElementById(`comments-${postId}`);
//             const commentHTML = `
//                 <div class="comment">
//                     <div class="comment-meta">
//                         <span class="comment-author">${comment.Author}</span>
//                         <span class="comment-date">${formatDate(comment.CreatedAt)}</span>
//                     </div>
//                     <p class="comment-content">${comment.Content}</p>
//                     <div class="comment-actions">
//                         <button class="like-btn" data-comment-id="${comment.ID}">
//                             👍 <span class="like-count">${comment.LikeCount}</span>
//                         </button>
//                         <button class="dislike-btn" data-comment-id="${comment.ID}">
//                             👎 <span class="dislike-count">${comment.DislikeCount}</span>
//                         </button>
//                     </div>
//                 </div>
//             `;
//             commentsSection.insertAdjacentHTML('afterbegin', commentHTML);
            
//             // Clear the form
//             form.querySelector('textarea').value = '';
//         } catch (error) {
//             console.error('Error creating comment:', error);
//             alert('Failed to create comment. Please try again.');
//         }
//     }
// });

// Helper function to format dates
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
        month: 'short',
        day: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}
