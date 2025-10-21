// TinyGo Web UI JavaScript
class TinyGoApp {
    constructor() {
        this.baseURL = window.location.origin;
        this.init();
    }

    init() {
        this.bindEvents();
        this.loadStats();
    }

    bindEvents() {
        // 绑定表单提交事件
        const form = document.getElementById('shortenForm');
        if (form) {
            form.addEventListener('submit', (e) => this.handleShorten(e));
        }

        // 绑定删除按钮事件
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('delete-btn')) {
                this.handleDelete(e.target.dataset.code);
            }
        });
    }

    async handleShorten(e) {
        e.preventDefault();
        
        const form = e.target;
        const formData = new FormData(form);
        const longURL = formData.get('long_url');
        const customCode = formData.get('custom_code');

        if (!longURL) {
            this.showError('请输入要缩短的 URL');
            return;
        }

        try {
            this.showLoading(true);
            
            const response = await fetch(`${this.baseURL}/api/shorten`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    long_url: longURL,
                    custom_code: customCode || undefined
                })
            });

            const result = await response.json();

            if (response.ok) {
                this.showSuccess(`短链接创建成功！`);
                this.displayShortLink(result);
                form.reset();
                this.loadStats(); // 刷新统计信息
            } else {
                this.showError(result.error || '创建短链接失败');
            }
        } catch (error) {
            this.showError('网络错误：' + error.message);
        } finally {
            this.showLoading(false);
        }
    }

    async handleDelete(code) {
        if (!confirm('确定要删除这个短链接吗？')) {
            return;
        }

        try {
            const response = await fetch(`${this.baseURL}/api/links/${code}`, {
                method: 'DELETE'
            });

            if (response.ok) {
                this.showSuccess('短链接已删除');
                this.loadStats(); // 刷新统计信息
            } else {
                const result = await response.json();
                this.showError(result.error || '删除失败');
            }
        } catch (error) {
            this.showError('网络错误：' + error.message);
        }
    }

    async loadStats() {
        try {
            const response = await fetch(`${this.baseURL}/admin/stats`);
            const stats = await response.json();

            if (response.ok) {
                this.updateStats(stats);
                this.updateLinksList(stats.links);
            }
        } catch (error) {
            console.error('加载统计信息失败：', error);
        }
    }

    updateStats(stats) {
        const totalLinksEl = document.getElementById('totalLinks');
        const totalHitsEl = document.getElementById('totalHits');
        
        if (totalLinksEl) totalLinksEl.textContent = stats.total_links;
        if (totalHitsEl) totalHitsEl.textContent = stats.total_hits;
    }

    updateLinksList(links) {
        const tbody = document.querySelector('#linksTable tbody');
        if (!tbody) return;

        tbody.innerHTML = '';

        if (links.length === 0) {
            tbody.innerHTML = '<tr><td colspan="5" class="text-center">暂无短链接</td></tr>';
            return;
        }

        links.forEach(link => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>
                    <a href="${this.baseURL}/${link.code}" target="_blank" class="short-link">
                        ${this.baseURL}/${link.code}
                    </a>
                </td>
                <td>
                    <a href="${link.long_url}" target="_blank" class="long-link" title="${link.long_url}">
                        ${this.truncateURL(link.long_url, 50)}
                    </a>
                </td>
                <td>${link.hit_count}</td>
                <td>${this.formatDate(link.created_at)}</td>
                <td>
                    <button class="btn btn-danger btn-sm delete-btn" data-code="${link.code}">
                        删除
                    </button>
                </td>
            `;
            tbody.appendChild(row);
        });
    }

    displayShortLink(result) {
        const resultDiv = document.getElementById('result');
        if (!resultDiv) return;

        resultDiv.innerHTML = `
            <div class="alert alert-success">
                <h4>短链接创建成功！</h4>
                <p><strong>短链接：</strong> 
                    <a href="${result.short_url}" target="_blank" class="short-link">
                        ${result.short_url}
                    </a>
                </p>
                <p><strong>原始链接：</strong> 
                    <a href="${result.long_url}" target="_blank" class="long-link">
                        ${result.long_url}
                    </a>
                </p>
                <button class="btn btn-primary" onclick="navigator.clipboard.writeText('${result.short_url}')">
                    复制短链接
                </button>
            </div>
        `;
        resultDiv.style.display = 'block';
    }

    showSuccess(message) {
        this.showMessage(message, 'success');
    }

    showError(message) {
        this.showMessage(message, 'error');
    }

    showMessage(message, type) {
        const alertClass = type === 'success' ? 'alert-success' : 'alert-danger';
        const messageDiv = document.createElement('div');
        messageDiv.className = `alert ${alertClass} alert-dismissible fade show`;
        messageDiv.innerHTML = `
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        `;

        const container = document.querySelector('.container');
        if (container) {
            container.insertBefore(messageDiv, container.firstChild);
            
            // 3秒后自动消失
            setTimeout(() => {
                if (messageDiv.parentNode) {
                    messageDiv.remove();
                }
            }, 3000);
        }
    }

    showLoading(show) {
        const submitBtn = document.querySelector('#shortenForm button[type="submit"]');
        if (submitBtn) {
            submitBtn.disabled = show;
            submitBtn.textContent = show ? '创建中...' : '创建短链接';
        }
    }

    truncateURL(url, maxLength) {
        if (url.length <= maxLength) return url;
        return url.substring(0, maxLength) + '...';
    }

    formatDate(dateString) {
        const date = new Date(dateString);
        return date.toLocaleString('zh-CN');
    }
}

// 页面加载完成后初始化应用
document.addEventListener('DOMContentLoaded', () => {
    new TinyGoApp();
});
