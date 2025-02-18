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
        
            window.location.reload();
        } catch (error) {
            alert('Failed to create comment. Please try again.');
        }
    }
});
