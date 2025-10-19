// Telegram Web App initialization
const tg = window.Telegram.WebApp;

// Initialize the app
tg.ready();
tg.expand();

// Get user data from Telegram
const user = tg.initDataUnsafe.user;
let currentUser = null;

// API base URL
const API_BASE = '/api';

// Добавляем заголовок для обхода предупреждения ngrok
const originalFetch = window.fetch;
window.fetch = function(url, options = {}) {
    options.headers = {
        ...options.headers,
        'ngrok-skip-browser-warning': 'true'
    };
    return originalFetch(url, options);
};

// Initialize the app
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
});

async function initializeApp() {
    try {
        // Update user info in header
        if (user) {
            document.getElementById('user-name').textContent = 
                `${user.first_name} ${user.last_name || ''}`.trim();
            
            // Create or get user in our database
            await createOrGetUser();
        }
        
        // Load initial data
        await loadTeams();
    } catch (error) {
        console.error('Error initializing app:', error);
        showError('Ошибка инициализации приложения');
    }
}

async function createOrGetUser() {
    if (!user) return;
    
    try {
        // Try to get existing user
        const response = await fetch(`${API_BASE}/users/${user.id}`);
        if (response.ok) {
            currentUser = await response.json();
        } else {
            // Create new user
            const userData = {
                telegram_id: user.id,
                username: user.username || '',
                first_name: user.first_name || '',
                last_name: user.last_name || ''
            };
            
            const createResponse = await fetch(`${API_BASE}/users`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData)
            });
            
            if (createResponse.ok) {
                currentUser = await createResponse.json();
            }
        }
    } catch (error) {
        console.error('Error creating/getting user:', error);
    }
}

async function loadTeams() {
    try {
        const response = await fetch(`${API_BASE}/teams`);
        if (response.ok) {
            const teams = await response.json();
            displayTeams(teams);
        } else {
            showError('Ошибка загрузки команд');
        }
    } catch (error) {
        console.error('Error loading teams:', error);
        showError('Ошибка загрузки команд');
    }
}

function displayTeams(teams) {
    const teamsList = document.getElementById('teams-list');
    
    if (teams.length === 0) {
        teamsList.innerHTML = '<div class="loading">Команды не найдены</div>';
        return;
    }
    
    teamsList.innerHTML = teams.map(team => `
        <div class="team-card">
            <div class="team-name">${team.name}</div>
            <div class="team-description">${team.description || 'Нет описания'}</div>
            <div class="team-meta">
                <span>Участников: ${team.members ? team.members.length : 0}</span>
                <div class="team-actions">
                    ${getTeamActionButton(team)}
                </div>
            </div>
        </div>
    `).join('');
}

function getTeamActionButton(team) {
    if (!currentUser) return '';
    
    const isInTeam = currentUser.team_id === team.id;
    const isMyTeam = currentUser.team_id === team.id;
    
    if (isInTeam) {
        return `<button class="btn btn-danger" onclick="leaveTeam()">Покинуть команду</button>`;
    } else {
        return `<button class="btn btn-success" onclick="joinTeam(${team.id})">Вступить в команду</button>`;
    }
}

async function joinTeam(teamId) {
    if (!currentUser) {
        showError('Пользователь не найден');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/teams/${teamId}/join/${currentUser.telegram_id}`, {
            method: 'POST'
        });
        
        if (response.ok) {
            showSuccess('Вы успешно вступили в команду!');
            await loadTeams();
            await loadMyTeam();
        } else {
            const error = await response.json();
            showError(error.error || 'Ошибка вступления в команду');
        }
    } catch (error) {
        console.error('Error joining team:', error);
        showError('Ошибка вступления в команду');
    }
}

async function leaveTeam() {
    if (!currentUser) {
        showError('Пользователь не найден');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/teams/leave/${currentUser.telegram_id}`, {
            method: 'POST'
        });
        
        if (response.ok) {
            showSuccess('Вы покинули команду');
            await loadTeams();
            await loadMyTeam();
        } else {
            const error = await response.json();
            showError(error.error || 'Ошибка выхода из команды');
        }
    } catch (error) {
        console.error('Error leaving team:', error);
        showError('Ошибка выхода из команды');
    }
}

async function loadMyTeam() {
    if (!currentUser || !currentUser.team_id) {
        document.getElementById('my-team-content').innerHTML = 
            '<div class="loading">Вы не состоите ни в одной команде</div>';
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/teams/${currentUser.team_id}`);
        if (response.ok) {
            const team = await response.json();
            displayMyTeam(team);
        } else {
            document.getElementById('my-team-content').innerHTML = 
                '<div class="loading">Ошибка загрузки команды</div>';
        }
    } catch (error) {
        console.error('Error loading my team:', error);
        document.getElementById('my-team-content').innerHTML = 
            '<div class="loading">Ошибка загрузки команды</div>';
    }
}

function displayMyTeam(team) {
    const myTeamContent = document.getElementById('my-team-content');
    
    myTeamContent.innerHTML = `
        <div class="team-card">
            <div class="team-name">${team.name}</div>
            <div class="team-description">${team.description || 'Нет описания'}</div>
            <div class="team-meta">
                <span>Участников: ${team.members ? team.members.length : 0}</span>
                <button class="btn btn-danger" onclick="leaveTeam()">Покинуть команду</button>
            </div>
            ${team.members && team.members.length > 0 ? `
                <div class="team-members">
                    <h4>Участники:</h4>
                    <ul>
                        ${team.members.map(member => `
                            <li>${member.first_name} ${member.last_name || ''} (@${member.username || 'no_username'})</li>
                        `).join('')}
                    </ul>
                </div>
            ` : ''}
        </div>
    `;
}

function showTab(tabName) {
    // Hide all tabs
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.classList.remove('active');
    });
    
    // Remove active class from all buttons
    document.querySelectorAll('.tab-button').forEach(button => {
        button.classList.remove('active');
    });
    
    // Show selected tab
    document.getElementById(tabName + '-tab').classList.add('active');
    
    // Add active class to clicked button
    event.target.classList.add('active');
    
    // Load content for the tab
    if (tabName === 'my-team') {
        loadMyTeam();
    }
}

function showCreateTeamModal() {
    document.getElementById('create-team-modal').style.display = 'block';
}

function closeCreateTeamModal() {
    document.getElementById('create-team-modal').style.display = 'none';
    document.getElementById('create-team-form').reset();
}

// Handle create team form submission
document.getElementById('create-team-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const teamData = {
        name: formData.get('name'),
        description: formData.get('description'),
        created_by: currentUser ? currentUser.id : null
    };
    
    try {
        const response = await fetch(`${API_BASE}/teams`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(teamData)
        });
        
        if (response.ok) {
            showSuccess('Команда успешно создана!');
            closeCreateTeamModal();
            await loadTeams();
        } else {
            const error = await response.json();
            showError(error.error || 'Ошибка создания команды');
        }
    } catch (error) {
        console.error('Error creating team:', error);
        showError('Ошибка создания команды');
    }
});

function showError(message) {
    // Remove existing messages
    document.querySelectorAll('.error, .success').forEach(el => el.remove());
    
    const errorDiv = document.createElement('div');
    errorDiv.className = 'error';
    errorDiv.textContent = message;
    
    document.querySelector('.container').insertBefore(errorDiv, document.querySelector('main'));
    
    // Auto remove after 5 seconds
    setTimeout(() => errorDiv.remove(), 5000);
}

function showSuccess(message) {
    // Remove existing messages
    document.querySelectorAll('.error, .success').forEach(el => el.remove());
    
    const successDiv = document.createElement('div');
    successDiv.className = 'success';
    successDiv.textContent = message;
    
    document.querySelector('.container').insertBefore(successDiv, document.querySelector('main'));
    
    // Auto remove after 5 seconds
    setTimeout(() => successDiv.remove(), 5000);
}

// Close modal when clicking outside
window.onclick = function(event) {
    const modal = document.getElementById('create-team-modal');
    if (event.target === modal) {
        closeCreateTeamModal();
    }
}
