// Global variables
const API_BASE_URL = 'http://localhost:8080/api/v1';
let currentUser = null;

// Tab functionality
function showTab(tabName) {
    // Hide all tab contents
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.classList.remove('active');
    });
    
    // Remove active class from all tab buttons
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    
    // Show selected tab
    document.getElementById(tabName + '-tab').classList.add('active');
    
    // Add active class to clicked button
    event.target.classList.add('active');
    
    // Load data for specific tabs
    if (tabName === 'manage') {
        loadEmergencyList();
    } else if (tabName === 'history') {
        loadEmergencyHistory();
    }
}

// Emergency form submission
document.getElementById('emergency-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const formData = new FormData();
    formData.append('user_id', document.getElementById('user_id').value);
    formData.append('type', document.getElementById('type').value);
    formData.append('title', document.getElementById('title').value);
    formData.append('description', document.getElementById('description').value);
    formData.append('location', document.getElementById('location').value);
    formData.append('map_link', document.getElementById('map_link').value);
    
    const fileInput = document.getElementById('file');
    if (fileInput.files[0]) {
        formData.append('file', fileInput.files[0]);
    }
    
    try {
        showLoading('btn-submit');
        
        const response = await fetch(`${API_BASE_URL}/emergencies/create`, {
            method: 'POST',
            body: formData
        });
        
        const result = await response.json();
        
        if (response.ok) {
            showAlert('รายงานเหตุฉุกเฉินสำเร็จ!', 'success');
            document.getElementById('emergency-form').reset();
            document.getElementById('file-preview').innerHTML = '';
        } else {
            showAlert(result.message || 'เกิดข้อผิดพลาด', 'error');
        }
    } catch (error) {
        console.error('Error:', error);
        showAlert('เกิดข้อผิดพลาดในการเชื่อมต่อ', 'error');
    } finally {
        hideLoading('btn-submit');
    }
});

// File preview functionality
document.getElementById('file').addEventListener('change', function(e) {
    const file = e.target.files[0];
    const preview = document.getElementById('file-preview');
    
    if (file) {
        if (file.type.startsWith('image/')) {
            const reader = new FileReader();
            reader.onload = function(e) {
                preview.innerHTML = `
                    <div style="text-align: center;">
                        <img src="${e.target.result}" style="max-width: 200px; max-height: 200px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
                        <p style="margin-top: 10px; color: #27ae60; font-weight: 600;">
                            <i class="fas fa-check-circle"></i> ${file.name}
                        </p>
                    </div>
                `;
            };
            reader.readAsDataURL(file);
        } else {
            preview.innerHTML = `
                <div style="color: #e74c3c;">
                    <i class="fas fa-exclamation-triangle"></i> 
                    กรุณาเลือกไฟล์รูปภาพเท่านั้น
                </div>
            `;
        }
    } else {
        preview.innerHTML = '';
    }
});

// Update emergency status (Officer)
async function updateEmergency(emergencyId, status, actionNote = '') {
    const formData = new FormData();
    formData.append('officer_id', 'off-1234-5678'); // This should come from logged in officer
    formData.append('status', status);
    if (actionNote) {
        formData.append('action_note', actionNote);
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/emergencies/officer/${emergencyId}`, {
            method: 'PATCH',
            body: formData
        });
        
        const result = await response.json();
        
        if (response.ok) {
            showAlert(`อัปเดตสถานะเป็น "${status}" สำเร็จ!`, 'success');
            loadEmergencyList(); // Refresh the list
        } else {
            showAlert(result.message || 'เกิดข้อผิดพลาดในการอัปเดต', 'error');
        }
    } catch (error) {
        console.error('Error:', error);
        showAlert('เกิดข้อผิดพลาดในการเชื่อมต่อ', 'error');
    }
}

// Load emergency list for management
async function loadEmergencyList() {
    try {
        const response = await fetch(`${API_BASE_URL}/emergencies/list`);
        const result = await response.json();
        
        if (response.ok && result.data) {
            renderEmergencyList(result.data);
        } else {
            console.error('Failed to load emergency list:', result);
        }
    } catch (error) {
        console.error('Error loading emergency list:', error);
    }
}

// Render emergency list
function renderEmergencyList(emergencies) {
    const container = document.querySelector('.emergency-list');
    
    if (!emergencies || emergencies.length === 0) {
        container.innerHTML = `
            <div style="text-align: center; padding: 40px; color: #7f8c8d;">
                <i class="fas fa-inbox" style="font-size: 3rem; margin-bottom: 20px;"></i>
                <h3>ไม่มีเหตุฉุกเฉินในขณะนี้</h3>
                <p>ขณะนี้ไม่มีเหตุฉุกเฉินที่ต้องจัดการ</p>
            </div>
        `;
        return;
    }
    
    container.innerHTML = emergencies.map(emergency => {
        const statusClass = getStatusClass(emergency.status);
        const typeIcon = getTypeIcon(emergency.type);
        const timeAgo = getTimeAgo(emergency.created_at);
        
        return `
            <div class="emergency-card ${statusClass}">
                <div class="card-header">
                    <span class="emergency-type">${typeIcon} ${emergency.type}</span>
                    <span class="status-badge ${statusClass}">${emergency.status}</span>
                </div>
                <h3>${emergency.title}</h3>
                <p class="description">${emergency.description}</p>
                <div class="emergency-details">
                    <span><i class="fas fa-map-marker-alt"></i> ${emergency.location}</span>
                    <span><i class="fas fa-clock"></i> ${timeAgo}</span>
                    ${emergency.officer_id ? `<span><i class="fas fa-user"></i> เจ้าหน้าที่: ${emergency.officer_id}</span>` : ''}
                </div>
                <div class="action-buttons">
                    ${emergency.status === 'รอการตอบสนอง' ? `
                        <button class="btn-action accept" onclick="updateEmergency('${emergency.id}', 'กำลังดำเนินการ')">
                            <i class="fas fa-check"></i> รับเรื่อง
                        </button>
                    ` : ''}
                    <button class="btn-action view" onclick="viewDetails('${emergency.id}')">
                        <i class="fas fa-eye"></i> ดูรายละเอียด
                    </button>
                    ${emergency.status === 'กำลังดำเนินการ' ? `
                        <button class="btn-action complete" onclick="updateEmergency('${emergency.id}', 'เสร็จสิ้น', 'ดำเนินการเสร็จสิ้น')">
                            <i class="fas fa-check-circle"></i> เสร็จสิ้น
                        </button>
                    ` : ''}
                </div>
                ${emergency.status === 'กำลังดำเนินการ' ? `
                    <div class="action-form">
                        <textarea id="action-note-${emergency.id}" placeholder="บันทึกการดำเนินการ..." rows="3"></textarea>
                        <div class="form-actions">
                            <button class="btn-action update" onclick="updateEmergencyWithNote('${emergency.id}')">
                                <i class="fas fa-edit"></i> อัปเดต
                            </button>
                        </div>
                    </div>
                ` : ''}
            </div>
        `;
    }).join('');
}

// Update emergency with action note
function updateEmergencyWithNote(emergencyId) {
    const actionNote = document.getElementById(`action-note-${emergencyId}`).value;
    if (actionNote.trim()) {
        updateEmergency(emergencyId, 'กำลังดำเนินการ', actionNote);
        document.getElementById(`action-note-${emergencyId}`).value = '';
    } else {
        showAlert('กรุณาใส่บันทึกการดำเนินการ', 'error');
    }
}

// Load emergency history
async function loadEmergencyHistory() {
    try {
        const response = await fetch(`${API_BASE_URL}/emergencies/list`);
        const result = await response.json();
        
        if (response.ok && result.data) {
            renderEmergencyHistory(result.data.filter(e => e.status === 'เสร็จสิ้น'));
        }
    } catch (error) {
        console.error('Error loading emergency history:', error);
    }
}

// Render emergency history
function renderEmergencyHistory(emergencies) {
    const container = document.querySelector('.history-list');
    
    if (!emergencies || emergencies.length === 0) {
        container.innerHTML = `
            <div style="text-align: center; padding: 40px; color: #7f8c8d;">
                <i class="fas fa-history" style="font-size: 3rem; margin-bottom: 20px;"></i>
                <h3>ไม่มีประวัติ</h3>
                <p>ยังไม่มีเหตุฉุกเฉินที่เสร็จสิ้นแล้ว</p>
            </div>
        `;
        return;
    }
    
    container.innerHTML = emergencies.map(emergency => {
        const typeIcon = getTypeIcon(emergency.type);
        const timeAgo = getTimeAgo(emergency.updated_at || emergency.created_at);
        
        return `
            <div class="history-item completed">
                <div class="item-info">
                    <span class="emergency-type">${typeIcon} ${emergency.type}</span>
                    <h4>${emergency.title}</h4>
                    <p>${emergency.location}</p>
                </div>
                <div class="item-status">
                    <span class="status-badge completed">เสร็จสิ้น</span>
                    <span class="date">${timeAgo}</span>
                </div>
            </div>
        `;
    }).join('');
}

// View emergency details
function viewDetails(emergencyId) {
    // This would typically fetch detailed data
    document.getElementById('emergency-modal').style.display = 'block';
}

// Close modal
function closeModal() {
    document.getElementById('emergency-modal').style.display = 'none';
}

// Close modal when clicking outside
window.onclick = function(event) {
    const modal = document.getElementById('emergency-modal');
    if (event.target === modal) {
        modal.style.display = 'none';
    }
}

// Utility functions
function getStatusClass(status) {
    switch (status) {
        case 'รอการตอบสนอง': return 'pending';
        case 'กำลังดำเนินการ': return 'in-progress';
        case 'เสร็จสิ้น': return 'completed';
        default: return 'pending';
    }
}

function getTypeIcon(type) {
    const icons = {
        'ไฟไหม้': '🔥',
        'อุบัติเหตุ': '🚗',
        'ทะเลาะวิวาท': '⚔️',
        'ลักขโมย': '🔓',
        'แก๊สรั่ว': '💨',
        'น้ำท่วม': '🌊',
        'อื่นๆ': '📋'
    };
    return icons[type] || '📋';
}

function getTimeAgo(timestamp) {
    if (!timestamp) return 'ไม่ทราบเวลา';
    
    const now = Date.now();
    const time = timestamp * 1000; // Convert from unix timestamp
    const diff = now - time;
    
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);
    
    if (minutes < 1) return 'เมื่อสักครู่';
    if (minutes < 60) return `${minutes} นาทีที่แล้ว`;
    if (hours < 24) return `${hours} ชั่วโมงที่แล้ว`;
    return `${days} วันที่แล้ว`;
}

function showAlert(message, type) {
    // Remove existing alerts
    document.querySelectorAll('.alert').forEach(alert => alert.remove());
    
    const alert = document.createElement('div');
    alert.className = `alert ${type}`;
    alert.innerHTML = `
        <i class="fas fa-${type === 'success' ? 'check-circle' : type === 'error' ? 'exclamation-circle' : 'info-circle'}"></i>
        ${message}
    `;
    
    // Insert at the top of the active tab content
    const activeTab = document.querySelector('.tab-content.active');
    activeTab.insertBefore(alert, activeTab.firstChild);
    
    // Auto remove after 5 seconds
    setTimeout(() => {
        alert.remove();
    }, 5000);
}

function showLoading(buttonId) {
    const button = document.getElementById(buttonId) || document.querySelector(`.${buttonId}`);
    if (button) {
        button.disabled = true;
        button.innerHTML = '<span class="loading"></span> กำลังส่ง...';
    }
}

function hideLoading(buttonId) {
    const button = document.getElementById(buttonId) || document.querySelector(`.${buttonId}`);
    if (button) {
        button.disabled = false;
        button.innerHTML = '<i class="fas fa-paper-plane"></i> ส่งรายงาน';
    }
}

// Initialize page
document.addEventListener('DOMContentLoaded', function() {
    // Load emergency list on page load
    loadEmergencyList();
    
    // Add event listeners for filters
    document.querySelectorAll('.filter-select, .search-input').forEach(element => {
        element.addEventListener('change', function() {
            // Implement filtering logic here
            console.log('Filter changed:', this.value);
        });
    });
});

// Real-time updates (you could implement WebSocket here)
setInterval(() => {
    const activeTab = document.querySelector('.tab-btn.active');
    if (activeTab && activeTab.textContent.includes('จัดการเหตุฉุกเฉิน')) {
        loadEmergencyList();
    }
}, 30000); // Refresh every 30 seconds
