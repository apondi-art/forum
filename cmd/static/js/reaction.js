// Like/Dislike handling
document.addEventListener('click', async (e) => {
    if (e.target.classList.contains('like-btn') || e.target.classList.contains('dislike-btn')) {
        const button = e.target;
        const targetType = button.dataset.targetType;  // 'post' or 'comment'
        const targetId = button.dataset.targetId  // ID of the post or comment
        const type = button.classList.contains('like-btn') ? 'like' : 'dislike';
        try {
            const response = await fetch('/api/reaction', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ 
                    targetId: Number(targetId),
                    targetType: targetType,
                    type: type,
                 })
            });
            if (!response.ok) {
                throw new Error('Failed to update reaction');
            }
            const data = await response.json();
            
            // Update like/dislike counts
            const postCard = button.closest('.post-card, .comment');
            postCard.querySelector('.like-count').textContent = data.likes;
            postCard.querySelector('.dislike-count').textContent = data.dislikes;

            // toggle active class
            button.classList.toggle('active')
        } catch (error) {
            console.error('Error updating reaction:', error);
            alert('Failed to update reaction. Please try again.');
        }
    }
});

// Add CSS for active state
const style = document.createElement('style');
style.textContent = `
    .like-btn.active, .dislike-btn.active {
        background-color: #D34500;
        color: white;
    }
`;
document.head.appendChild(style);