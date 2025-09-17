document.addEventListener('DOMContentLoaded', function() {
    // Tab switching
    window.switchTab = function(tabName) {
        // Remove active class from all tabs and content
        document.querySelectorAll('.tab').forEach(tab => tab.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));

        // Add active class to selected tab and content
        document.querySelector(`[onclick="switchTab('${tabName}')"]`).classList.add('active');
        document.getElementById(`${tabName}-tab`).classList.add('active');

        // Update hidden input
        document.getElementById('upload_method').value = tabName;
    };

    // File input handling
    const fileInput = document.getElementById('file');
    const fileInfo = document.getElementById('file-info');
    const uploadArea = document.querySelector('.upload-area');

    fileInput.addEventListener('change', function(e) {
        const file = e.target.files[0];
        if (file) {
            fileInfo.style.display = 'block';
            fileInfo.innerHTML = `
                <strong>Selected file:</strong> ${file.name}<br>
                <strong>Size:</strong> ${formatFileSize(file.size)}<br>
                <strong>Type:</strong> ${file.type || 'Unknown'}
            `;
        } else {
            fileInfo.style.display = 'none';
        }
    });

    // Drag and drop functionality
    uploadArea.addEventListener('dragover', function(e) {
        e.preventDefault();
        uploadArea.classList.add('dragover');
    });

    uploadArea.addEventListener('dragleave', function(e) {
        e.preventDefault();
        uploadArea.classList.remove('dragover');
    });

    uploadArea.addEventListener('drop', function(e) {
        e.preventDefault();
        uploadArea.classList.remove('dragover');

        const files = e.dataTransfer.files;
        if (files.length > 0) {
            fileInput.files = files;

            // Trigger change event
            const event = new Event('change');
            fileInput.dispatchEvent(event);
        }
    });

    // Form validation
    document.getElementById('uploadForm').addEventListener('submit', function(e) {
        const uploadMethod = document.getElementById('upload_method').value;

        if (uploadMethod === 'file') {
            const file = fileInput.files[0];
            if (!file) {
                e.preventDefault();
                alert('Please select a file to upload.');
                return false;
            }

            // Check file size (50MB)
            if (file.size > 50 * 1024 * 1024) {
                e.preventDefault();
                alert('File size must be less than 50MB.');
                return false;
            }

            // Check file type
            const validTypes = ['.yaml', '.yml', '.json'];
            const fileName = file.name.toLowerCase();
            const isValidType = validTypes.some(ext => fileName.endsWith(ext));
            if (!isValidType) {
                e.preventDefault();
                alert('Please upload a valid YAML or JSON file.');
                return false;
            }
        } else if (uploadMethod === 'paste') {
            const yamlContent = document.getElementById('yaml_content').value.trim();
            if (!yamlContent) {
                e.preventDefault();
                alert('Please paste your OpenAPI content.');
                return false;
            }
        }

        // Show loading state
        const submitBtn = this.querySelector('button[type="submit"]');
        submitBtn.disabled = true;
        submitBtn.innerHTML = '<span class="loading"></span> Processing...';
    });

    // Utility function to format file size
    function formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }
});