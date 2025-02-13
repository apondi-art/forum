// Comment section toggle
document.addEventListener('click', (e) => {
    if (e.target.classList.contains('view-comments')) {
        const postId = e.target.dataset.postId;
        const commentsSection = document.getElementById(`comments-${postId}`);
        commentsSection.style.display = commentsSection.style.display === 'none' ? 'block' : 'none';
    }
});

// Comment form handling
document.addEventListener('submit', async (e) => {
    if (e.target.classList.contains('comment-form')) {
        e.preventDefault();

        const form = e.target;
        const postId = form.dataset.postId;
        const content = form.querySelector('textarea').value;

        if (!postId || !content.trim()) {
            alert('Comment content cannot be empty.');
            return;
        }
        
        try {
            const response = await fetch("/api/comment", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ postId: Number(postId), content: content })
            });
             
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`Failed to create comment: ${errorText}`);
            }
        
            const comment = await response.json();
            
            // Add new comment to the page
            const commentsSection = document.getElementById(`comments-${postId}`);
            const commentHTML = `
                <div class="comment">
                    <div class="comment-meta">
                        <span class="comment-author">${comment.Author}</span>
                        <span class="comment-date">${formatDate(comment.CreatedAt)}</span>
                    </div>
                    <p class="comment-content">${comment.Content}</p>
                    <div class="comment-actions">
                        <button class="like-btn" data-comment-id="${comment.ID}">
                            👍 <span class="like-count">${comment.LikeCount}</span>
                        </button>
                        <button class="dislike-btn" data-comment-id="${comment.ID}">
                            👎 <span class="dislike-count">${comment.DislikeCount}</span>
                        </button>
                    </div>
                </div>
            `;
            commentsSection.insertAdjacentHTML('afterbegin', commentHTML);
            // Clear the form
            form.querySelector('textarea').value = '';
        } catch (error) {
            alert('Failed to create comment. Please try again.');
        }
    }
});

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