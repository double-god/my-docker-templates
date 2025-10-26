import React, { useState } from 'react';
import uploadToOSSWithSTS from './ossUploader'; // 导入函数
import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api';

const MyComponent = () => {
    const [uploading, setUploading] = useState(false);
    const [imageUrl, setImageUrl] = useState(''); // 保存上传后的 URL

    const handleFileChange = async (event) => {
        const file = event.target.files[0];
        if (file) {
            setUploading(true);
            setImageUrl(''); // 清空旧 URL
            try {
                const uploadedUrl = await uploadToOSSWithSTS(file);
                setImageUrl(uploadedUrl); // 保存 URL 供后续使用
                alert('图片上传成功!');
                // 你可以在这里就把 URL 发给后端，或者等用户提交整个表单时再发
            } catch (error) {
                alert(`上传失败: ${error.message || '未知错误'}`);
            } finally {
                setUploading(false);
            }
        }
    };

    // 假设有一个提交按钮，把 imageUrl 和其他数据一起发给后端
    const handleSubmit = async () => {
        if (!imageUrl) {
            alert("请先上传图片");
            return;
        }
        try {
            // 调用你的后端 API 保存数据 (这里只是示例)
            await axios.post(`${API_BASE_URL}/hitokoto`, {
                text: "示例文字",
                from: "示例来源",
                imageURL: imageUrl // 把上传后的 URL 发给后端
            });
            alert("数据提交成功!");
            setImageUrl(''); // 清空
        } catch (error) {
            alert("提交失败!");
        }
    };

    return (
        <div>
            <input
                type="file"
                onChange={handleFileChange}
                disabled={uploading}
                accept="image/*"
            />
            {uploading && <p>正在上传...</p>}
            {imageUrl && (
                <div>
                    <p>图片预览:</p>
                    <img src={imageUrl} alt="上传预览" width="200" />
                    {/* 注意：如果 Bucket 是私有的，这里直接显示会失败，需要后端提供签名 URL */}
                    <button onClick={handleSubmit}>提交数据</button>
                </div>
            )}
        </div>
    );
};

export default MyComponent;