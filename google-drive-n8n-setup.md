# Google Drive Credentials Setup for n8n

## 🎯 Goal: Set up Google Drive credentials to access resume files

### Step 1: Create Google Cloud Project

1. **Go to [Google Cloud Console](https://console.cloud.google.com/)**
2. **Sign in** with your Google account
3. **Create a new project** or select existing one:
   - Click "Select a project" → "New Project"
   - Name: `Resource Search Integration`
   - Click "Create"

### Step 2: Enable Google Drive API

1. **In your project, go to "APIs & Services" → "Library"**
2. **Search for "Google Drive API"**
3. **Click on "Google Drive API"**
4. **Click "Enable"**

### Step 3: Create Service Account (Recommended Method)

#### Option A: Service Account (Easiest)

1. **Go to "IAM & Admin" → "Service Accounts"**
2. **Click "Create Service Account"**
3. **Fill in details:**
   - **Service account name**: `n8n-drive-access`
   - **Description**: `Service account for n8n Google Drive access`
4. **Click "Create and Continue"**
5. **Skip role assignment** (click "Continue")
6. **Click "Done"**

#### Option B: OAuth2 (Alternative)

1. **Go to "APIs & Services" → "Credentials"**
2. **Click "Create Credentials" → "OAuth client ID"**
3. **Choose "Web application"**
4. **Add authorized redirect URIs:**
   - `http://localhost:5678/rest/oauth2-credential/callback`
   - `https://your-n8n-instance.com/rest/oauth2-credential/callback`

### Step 4: Download Credentials

#### For Service Account:
1. **Click on your service account**
2. **Go to "Keys" tab**
3. **Click "Add Key" → "Create new key"**
4. **Choose "JSON"**
5. **Download the JSON file**
6. **Save as `google-drive-credentials.json`**

#### For OAuth2:
1. **Download the JSON file from OAuth client**
2. **Save as `google-drive-credentials.json`**

### Step 5: Configure Folder Access

1. **Open Google Drive**
2. **Navigate to your resume folder**
3. **Right-click on the folder** → **Share**
4. **Add the service account email** (from the JSON file)
5. **Set permission to "Editor"**
6. **Click "Send"**

### Step 6: Get Folder ID

1. **Open your resume folder in Google Drive**
2. **Copy the folder ID from the URL:**
   ```
   https://drive.google.com/drive/folders/FOLDER_ID_HERE
   ```
3. **Save this folder ID** - you'll need it for the workflow

### Step 7: Set Up n8n Credentials

#### Method 1: Service Account (Recommended)

1. **In n8n, go to "Credentials"**
2. **Click "Create New"**
3. **Search for "Google Drive"**
4. **Select "Google Drive OAuth2 API"**
5. **Choose "Service Account"**
6. **Upload the JSON file** you downloaded
7. **Click "Save"**

#### Method 2: OAuth2

1. **In n8n, go to "Credentials"**
2. **Click "Create New"**
3. **Search for "Google Drive"**
4. **Select "Google Drive OAuth2 API"**
5. **Choose "OAuth2"**
6. **Fill in:**
   - **Client ID**: From your JSON file
   - **Client Secret**: From your JSON file
7. **Click "Save"**

### Step 8: Test the Connection

1. **Go to your workflow**
2. **Click on the "Search Resume Files in Drive" node**
3. **Select your Google Drive credential**
4. **Set the folder ID** (from Step 6)
5. **Test the connection**

## 🔧 Troubleshooting

### Issue: "Access denied" or "Permission denied"
**Solution:**
- Make sure the service account has access to the folder
- Check that the folder is shared with the service account email
- Verify the service account has "Editor" permission

### Issue: "Folder not found"
**Solution:**
- Double-check the folder ID is correct
- Make sure the folder is shared with the service account
- Try using the full folder URL

### Issue: "API not enabled"
**Solution:**
- Go back to Google Cloud Console
- Enable Google Drive API
- Wait a few minutes for changes to propagate

### Issue: "Invalid credentials"
**Solution:**
- Re-download the JSON file
- Make sure you're using the correct credential type
- Check that the JSON file is valid

## 📋 Quick Checklist

- [ ] Google Cloud project created
- [ ] Google Drive API enabled
- [ ] Service account created
- [ ] JSON credentials downloaded
- [ ] Folder shared with service account
- [ ] Folder ID copied
- [ ] n8n credential created
- [ ] Connection tested

## 🚀 Alternative: Use OAuth2 with Your Google Account

If you prefer to use your personal Google account:

1. **Use OAuth2 method** instead of service account
2. **Authorize with your Google account**
3. **Grant permissions** when prompted
4. **Use your personal Google Drive access**

## 📁 Folder Structure Example

Your Google Drive folder should contain:
```
Resumes/
├── john-doe-resume.pdf
├── jane-smith-cv.docx
├── mike-johnson-resume.pdf
└── sarah-wilson-cv.pdf
```

## 🔒 Security Notes

- **Keep the JSON file secure**
- **Don't commit it to version control**
- **Use environment variables** for production
- **Rotate credentials regularly**

## 📞 Need Help?

If you're still having issues:
1. **Check the exact error message**
2. **Verify folder permissions**
3. **Test with a simple file first**
4. **Check n8n execution logs**

The service account method is recommended as it's more reliable and doesn't require user interaction!
