// Global app utilities and HTMX configuration

// Show/hide loading overlay
function showLoading() {
    document.getElementById('loading').style.display = 'flex';
}

function hideLoading() {
    document.getElementById('loading').style.display = 'none';
}

// Toast notification system
function showToast(message, type = 'info', duration = 5000) {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    toast.innerHTML = message.replace(/\n/g, '<br>');
    
    container.appendChild(toast);
    
    // Auto-remove toast
    setTimeout(() => {
        toast.style.animation = 'slideOut 0.3s ease-out forwards';
        setTimeout(() => {
            container.removeChild(toast);
        }, 300);
    }, duration);
}

// Add slideOut animation
const style = document.createElement('style');
style.textContent = `
    @keyframes slideOut {
        from {
            transform: translateX(0);
            opacity: 1;
        }
        to {
            transform: translateX(100%);
            opacity: 0;
        }
    }
`;
document.head.appendChild(style);

// HTMX Global Configuration
document.addEventListener('DOMContentLoaded', function() {
    // Configure HTMX
    htmx.config.globalViewTransitions = true;
    
    // Global HTMX event listeners
    document.body.addEventListener('htmx:beforeRequest', function(evt) {
        // Show loading for non-background requests
        if (!evt.detail.target.hasAttribute('hx-swap-oob')) {
            showLoading();
        }
    });

    document.body.addEventListener('htmx:afterRequest', function(evt) {
        hideLoading();
        
        // Handle errors globally
        if (evt.detail.xhr.status >= 400) {
            try {
                const response = JSON.parse(evt.detail.xhr.responseText);
                showToast(response.error || 'Erro no servidor', 'error');
            } catch (e) {
                showToast('Erro de conexão. Verifique sua internet.', 'error');
            }
        }
    });

    document.body.addEventListener('htmx:responseError', function(evt) {
        hideLoading();
        showToast('Erro de conexão. Tente novamente.', 'error');
    });

    document.body.addEventListener('htmx:timeout', function(evt) {
        hideLoading();
        showToast('Timeout na requisição. Tente novamente.', 'error');
    });

    // Set HTMX timeout to 10 seconds
    htmx.config.timeout = 10000;
});

// Form utilities
function serializeFormToJSON(form) {
    const formData = new FormData(form);
    const object = {};
    formData.forEach((value, key) => {
        if (object[key]) {
            if (!Array.isArray(object[key])) {
                object[key] = [object[key]];
            }
            object[key].push(value);
        } else {
            object[key] = value;
        }
    });
    return object;
}

// Date/time utilities
function formatDateTime(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString('pt-BR', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}

function isDateInFuture(dateString) {
    return new Date(dateString) > new Date();
}

function getTimeRemaining(dateString) {
    const now = new Date();
    const target = new Date(dateString);
    const diff = target - now;
    
    if (diff <= 0) return null;
    
    const hours = Math.floor(diff / (1000 * 60 * 60));
    const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
    
    if (hours > 0) {
        return `${hours}h ${minutes}m`;
    } else {
        return `${minutes}m`;
    }
}

// Question utilities
function renderQuestion(question) {
    const isAvailable = !isDateInFuture(question.scheduled_at);
    const timeRemaining = getTimeRemaining(question.scheduled_at);
    
    let statusClass = 'status-pending';
    let statusText = 'Aguardando';
    
    if (question.is_answered) {
        statusClass = 'status-answered';
        statusText = 'Respondida';
    } else if (isAvailable) {
        statusClass = 'status-available';
        statusText = 'Disponível';
    } else if (timeRemaining) {
        statusText = `Em ${timeRemaining}`;
    }
    
    return `
        <div class="question-card" data-question-id="${question.id}">
            <div class="question-header">
                <h3 class="question-title">${question.title}</h3>
                <div class="question-meta">
                    <span class="question-time">${formatDateTime(question.scheduled_at)}</span>
                    <span class="question-status ${statusClass}">${statusText}</span>
                </div>
            </div>
            <div class="question-content">
                ${question.description ? `<p class="question-description">${question.description}</p>` : ''}
                
                ${isAvailable && !question.is_answered ? `
                    <form class="question-form" 
                          hx-post="/api/answer" 
                          hx-trigger="submit"
                          hx-target="#answer-result-${question.id}"
                          hx-swap="innerHTML">
                        <input type="hidden" name="question_id" value="${question.id}">
                        <div class="question-options">
                            ${question.options.map(option => `
                                <button type="submit" name="answer" value="${option}" class="option-button">
                                    ${option}
                                </button>
                            `).join('')}
                        </div>
                    </form>
                    <div id="answer-result-${question.id}"></div>
                ` : ''}
                
                ${question.is_answered ? `
                    <div class="alert alert-success">
                        ✅ Pergunta já respondida!
                    </div>
                ` : ''}
                
                ${!isAvailable ? `
                    <div class="alert alert-info">
                        ⏰ Esta pergunta estará disponível em: ${timeRemaining || 'breve'}
                    </div>
                ` : ''}
                
                <div class="question-reward">
                    <div class="reward-label">🎁 Recompensa:</div>
                    <div class="reward-text">${question.reward}</div>
                </div>
            </div>
        </div>
    `;
}

function renderAdminQuestion(question) {
    const isActive = question.is_active;
    const isPast = new Date(question.scheduled_at) < new Date();
    
    return `
        <div class="question-card admin-question" data-question-id="${question.id}">
            <div class="question-header">
                <h3 class="question-title">${question.title}</h3>
                <div class="question-meta">
                    <span class="question-time">${formatDateTime(question.scheduled_at)}</span>
                    <span class="question-status ${isActive ? 'status-available' : 'status-pending'}">
                        ${isActive ? (isPast ? 'Ativa' : 'Agendada') : 'Inativa'}
                    </span>
                </div>
            </div>
            <div class="question-content">
                ${question.description ? `<p class="question-description">${question.description}</p>` : ''}
                
                <div class="question-options">
                    ${question.options.map((option, index) => `
                        <div class="option-preview ${option === question.correct_answer ? 'correct-answer' : ''}">
                            ${option} ${option === question.correct_answer ? '✅' : ''}
                        </div>
                    `).join('')}
                </div>
                
                <div class="question-reward">
                    <div class="reward-label">🎁 Recompensa:</div>
                    <div class="reward-text">${question.reward}</div>
                </div>
                
                <div class="admin-actions">
                    <button class="btn btn-secondary btn-small" onclick="editQuestion(${question.id})">
                        ✏️ Editar
                    </button>
                    <button class="btn btn-secondary btn-small" onclick="toggleQuestion(${question.id}, ${!isActive})">
                        ${isActive ? '⏸️ Desativar' : '▶️ Ativar'}
                    </button>
                </div>
            </div>
        </div>
    `;
}

// Dashboard rendering
function renderVisitorDashboard(data) {
    const stats = data.stats;
    const questions = data.questions || [];
    
    const statsHTML = `
        <div class="stat-card">
            <span class="stat-number">${stats.total_questions}</span>
            <span class="stat-label">Total de Perguntas</span>
        </div>
        <div class="stat-card">
            <span class="stat-number">${stats.available_questions}</span>
            <span class="stat-label">Disponíveis</span>
        </div>
        <div class="stat-card">
            <span class="stat-number">${stats.answered_questions}</span>
            <span class="stat-label">Respondidas</span>
        </div>
        <div class="stat-card">
            <span class="stat-number">${Math.round((stats.answered_questions / stats.total_questions) * 100) || 0}%</span>
            <span class="stat-label">Progresso</span>
        </div>
    `;
    
    document.getElementById('stats-container').innerHTML = statsHTML;
    
    const questionsHTML = questions.length > 0 
        ? `<div class="questions-grid">${questions.map(renderQuestion).join('')}</div>`
        : '<div class="alert alert-info">Nenhuma pergunta disponível ainda. Aguarde! 💕</div>';
    
    return questionsHTML;
}

function renderAdminDashboard(data) {
    const stats = data.stats;
    const questions = data.questions || [];
    
    const statsHTML = `
        <div class="stat-card">
            <span class="stat-number">${stats.total_questions}</span>
            <span class="stat-label">Total de Perguntas</span>
        </div>
        <div class="stat-card">
            <span class="stat-number">${stats.active_questions}</span>
            <span class="stat-label">Ativas</span>
        </div>
        <div class="stat-card">
            <span class="stat-number">${questions.filter(q => new Date(q.scheduled_at) > new Date()).length}</span>
            <span class="stat-label">Agendadas</span>
        </div>
        <div class="stat-card">
            <span class="stat-number">${questions.filter(q => new Date(q.scheduled_at) <= new Date()).length}</span>
            <span class="stat-label">Liberadas</span>
        </div>
    `;
    
    document.getElementById('stats-container').innerHTML = statsHTML;
    
    const questionsHTML = questions.length > 0 
        ? `<div class="questions-grid">${questions.map(renderAdminQuestion).join('')}</div>`
        : '<div class="alert alert-info">Nenhuma pergunta criada ainda. Crie a primeira! 🎯</div>';
    
    return questionsHTML;
}

// Auto-refresh utility
let autoRefreshInterval;

function startAutoRefresh(callback, interval = 30000) {
    stopAutoRefresh();
    autoRefreshInterval = setInterval(callback, interval);
}

function stopAutoRefresh() {
    if (autoRefreshInterval) {
        clearInterval(autoRefreshInterval);
        autoRefreshInterval = null;
    }
}

// Cleanup on page unload
window.addEventListener('beforeunload', stopAutoRefresh);

// Export utilities for global use
window.QuizApp = {
    showLoading,
    hideLoading,
    showToast,
    formatDateTime,
    isDateInFuture,
    getTimeRemaining,
    renderQuestion,
    renderAdminQuestion,
    renderVisitorDashboard,
    renderAdminDashboard,
    startAutoRefresh,
    stopAutoRefresh
};
