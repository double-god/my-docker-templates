import OSS from 'ali-oss';
import axios from 'axios';

// --- 读取 Vite 注入的环境变量 ---
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'; // 提供一个默认值

const uploadToOSSWithSTS = async (file) => {
    try {
        // 1. 从你的后端获取 STS 临时凭证 (使用环境变量配置的 API 地址)
        const response = await axios.get(`${API_BASE_URL}/oss/sts-credentials`); // <--- 修改点
        const creds = response.data;

        // 2. 初始化 OSS Client (代码不变)
        const client = new OSS({
            region: creds.region,
            accessKeyId: creds.accessKeyId,
            accessKeySecret: creds.accessKeySecret,
            stsToken: creds.stsToken,
            bucket: creds.bucket,
            secure: true,
        });

        // 3. 定义 Object Key (代码不变)
        const objectKey = `uploads/${Date.now()}-${file.name}`;

        // 4. 上传文件 (代码不变)
        console.log(`开始上传: ${objectKey}`);
        const result = await client.multipartUpload(objectKey, file, {
            progress: (p) => {
                const percent = Math.round(p * 100);
                console.log(`上传进度: ${percent}%`);
                // 更新 UI
            },
        });

        // 5. 获取 URL (代码不变)
        if (result && result.res && result.res.status === 200) {
            const fileURL = `https://${creds.bucket}.${creds.region}.aliyuncs.com/${objectKey}`;
            console.log('上传成功! 文件 URL:', fileURL);
            return fileURL;
        } else {
            throw new Error('OSS multipartUpload 返回结果异常');
        }

    } catch (error) {
        console.error('上传到 OSS (STS) 出错:', error);
        throw error;
    }
};

export default uploadToOSSWithSTS;