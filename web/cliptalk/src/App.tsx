import {Card, CardContent, CardDescription, CardHeader, CardTitle} from "@/components/ui/card"

import {Label} from "@/components/ui/label"
import {Input} from "@/components/ui/input"
import {Button} from "@/components/ui/button"
import React, {useState} from "react";
import {Checkbox} from "@/components/ui/checkbox";
import {ReloadIcon} from "@radix-ui/react-icons"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"

import {
    Drawer,
    DrawerClose,
    DrawerContent,
    DrawerDescription, DrawerFooter,
    DrawerHeader,
    DrawerTitle,
    DrawerTrigger,
} from "@/components/ui/drawer"

import useMediaQuery from "@/hooks/useMediaQuery.ts";
import SettingsIcon from "@/components/ui/Icon.tsx";

export default function App() {
    const [videoLink, setVideoLink] = useState("")
    const [isRemoveWatermark, setIsRemoveWatermark] = useState(false)
    const [videoFile, setVideoFile] = useState<File | null>(null)
    const [isLoading, setIsLoading] = useState(false);
    const [showDialog, setShowDialog] = useState(false);
    const [showSetting, setShowSetting] = useState(false);
    const isMobile = useMediaQuery("(max-width: 768px)");
    const [removeApiResponse, setRemoveApiResponse] = useState<{
        finalUrl: string,
        message: string,
        title: string
    } | null>(null);
    const [videoLinkApiResponse, setVideoLinkApiResponse] = useState<{
        content: string,
        duration: number,
        title: string,
        message: string,
        finalUrl: string
    } | null>(null);
    const [isGemini, setIsGmini] = useState(false)
    const [isOpenAI, setIsOpenAI] = useState(false)
    const [key,setKey] = useState('')

    // const[downloadLink, setDownloadLink] = useState("")
    // useEffect(() => {
    //     setIsRemoveWatermark(false)
    // }, []);
    function handleRemoveWatermarkChange() {
        setIsOpenAI(false)
        setIsGmini(false)
        setIsRemoveWatermark(!isRemoveWatermark)

    }

    function handleOpenAIChange() {
       setIsGmini(false)
        setIsRemoveWatermark(false)
        setIsOpenAI(!isOpenAI)
    }

    function handleGminiChange() {
        setIsOpenAI(false)
        setIsRemoveWatermark(false)
        setIsGmini(!isGemini)
    }

    function handleVideoLinkChange(event: React.ChangeEvent<HTMLInputElement>) {
        if (videoFile != null) {
            setVideoFile(null)
        }
        setVideoLink(event.target.value)
    }

    function handleVideoFileChange(event: React.ChangeEvent<HTMLInputElement>) {
        if (videoLink != '') {
            alert('请删除视频链接')
            event.target.value = ''; // 清除选定的文件
            return
        }
        const file = event.target.files?.[0];
        console.log('文件:', file);
        if (file) {
            const fileExtension = file.name.split('.').pop()?.toLowerCase();
            console.log('文件扩展名:', fileExtension);
            const videoExtensions = ['mp4', 'mkv', 'flv', 'avi', 'mov', 'wmv'];
            if (videoExtensions.includes(fileExtension || '')) {
                // 这是一个视频文件
                setVideoFile(file);
            } else {
                // 这不是一个视频文件
                alert('请上传一个视频文件');
                event.target.value = ''; // 清除选定的文件
            }
        }
    }

    async function makeRequest(url: string, data: any,key:string) {
        try {
            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    "Authorization": key,
                },
                body: JSON.stringify(data)
            });

            if (response.ok) {
                return await response.json();
            } else {
                console.error('Response:', response);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }

    async function mkFileRequest(url: string, data: any,key:string) {
        try {
            // 创建一个 FormData 对象
            const formData = new FormData();

            // 为 FormData 对象添加字段
            for (const key in data) {
                formData.append(key, data[key]);
            }

            const response = await fetch(url, {
                headers: {
                    "Authorization": key,
                },
                method: 'POST',
                body: formData, // 使用 FormData 对象作为请求体
            });

            if (response.ok) {
                return await response.json();
            } else {
                console.error('Response:', response);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }

    async function handleSubmit(event: React.FormEvent<HTMLButtonElement>) {
        event.preventDefault()

        if (videoFile == null && videoLink == '') {
            alert('请上传一个视频文件或者填写视频链接')
            return
        }
        if (videoFile != null && videoLink != '') {
            setVideoFile(null)
        }
        setIsLoading(true);
        if (videoLink != '' && isRemoveWatermark) {
            const requestBody = {
                url: videoLink
            };
            const data = await makeRequest('/remove', requestBody,"");
            if (data != null) {
                setRemoveApiResponse(data);
                setIsLoading(false);
                setShowDialog(true); // 显示对话框
            } else {
                alert('请求失败')
                setIsLoading(false);
            }
            return
        }
        let model = 'gemini'
        if (isOpenAI) {
            model = 'openai'
        } else if (isGemini) {
            model = 'gemini'
        }
        if(key==''){
            alert('请填写apikey')
            setShowSetting(true)
            setIsLoading(false)
            return
        }
        console.log(key)
        if (videoLink != '' && !isRemoveWatermark) {

            const requestBody = {
                url: videoLink,
                model: model
            };
            const data = await makeRequest('/video', requestBody,key);
            if (data != null) {
                setVideoLinkApiResponse(data);
                setIsLoading(false);
            } else {
                setIsLoading(false);
                alert('请求失败')
            }
            return
        }
        if (videoLink == '' && videoFile != null) {
            const requestBody = {
                file: videoFile,
                model: model
            };
            const data = await mkFileRequest('/video-file', requestBody,key);
            if (data != null) {
                setVideoLinkApiResponse(data);
                setIsLoading(false);
            } else {
                setIsLoading(false);
                alert('请求失败')
            }
            return
        }
    }

    function handleSettingSubmit(event: React.FormEvent<HTMLButtonElement>) {
        event.preventDefault()
        setShowSetting(true)
    }

    // const [open, setOpen] = useState(false)
    // const isDesktop = useMediaQuery("(min-width: 768px)")
    function apiKeySubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault()
        const formData = new FormData(e.currentTarget);

        // 使用 FormData.get 方法来获取表单中的数据
        const email = formData.getAll("key");
        // console.log('email:', email);
        setKey(email.toString())
        setShowSetting(false)

    }
    function DrawerDialogDemo() {
        if (!isMobile) {
            return (
                <Dialog open={showSetting} onOpenChange={setShowSetting}>
                    <DialogTrigger asChild>
                        {/*<Button variant="outline">Edit Profile</Button>*/}
                    </DialogTrigger>
                    <DialogContent className="sm:max-w-[425px]">
                        <DialogHeader>
                            <DialogTitle>Edit profile</DialogTitle>
                            <DialogDescription>
                                Make changes to your profile here. Click save when you're done.
                            </DialogDescription>
                        </DialogHeader>
                        <form className={"grid items-start gap-4"} onSubmit={apiKeySubmit}>
                            <div className="grid gap-2">
                                <Label htmlFor="key">apikey</Label>
                                <Input  name="key"/>
                            </div>
                            <Button type="submit" >Save changes</Button>
                        </form>
                    </DialogContent>
                </Dialog>
            )
        }

        return (
            <Drawer open={showSetting} onOpenChange={setShowSetting}>
                <DrawerTrigger asChild>
                    {/*<Button variant="outline">Edit Profile</Button>*/}
                </DrawerTrigger>
                <DrawerContent>
                    <DrawerHeader className="text-left">
                        <DrawerTitle>Edit profile</DrawerTitle>
                        <DrawerDescription>
                            Make changes to your profile here. Click save when you're done.
                        </DrawerDescription>
                    </DrawerHeader>
                    <form className={"grid items-start gap-4"} onSubmit={apiKeySubmit}>
                        <div className="grid gap-2">
                            <Label htmlFor="key">apikey</Label>
                            <Input  name="key"/>
                        </div>
                        <Button type="submit" >Save changes</Button>
                    </form>
                    <DrawerFooter className="pt-2">
                        <DrawerClose asChild>
                            <Button variant="outline">Cancel</Button>
                        </DrawerClose>
                    </DrawerFooter>
                </DrawerContent>
            </Drawer>
        )
    }

    function DownloadDialog() {
        return (
            isMobile ? (
                <Drawer open={showDialog} onOpenChange={setShowDialog}>
                    <DrawerTrigger asChild>
                        <Button variant="outline">Download</Button>
                    </DrawerTrigger>
                    <DrawerContent style={{height: '33.33vh'}}>
                        <DrawerHeader>
                            <DrawerTitle>Download Video</DrawerTitle>
                            <DrawerDescription>
                                Click the button below to download your video.
                            </DrawerDescription>
                            <h2>{[removeApiResponse?.title]}</h2>
                        </DrawerHeader>
                        <Button type="submit" onClick={() => window.open(removeApiResponse?.finalUrl)}>Download</Button>
                    </DrawerContent>
                </Drawer>
            ) : (
                <Dialog open={showDialog} onOpenChange={setShowDialog}>
                    <DialogTrigger asChild>
                        <Button variant="outline">Download</Button>
                    </DialogTrigger>
                    <DialogContent className="sm:max-w-[425px]">
                        <DialogHeader className="flex flex-col justify-center items-center">
                            <DialogTitle>Download Video</DialogTitle>
                            <DialogDescription>
                                Click the button below to download your video.
                            </DialogDescription>
                            <h2>{removeApiResponse?.title}</h2>
                        </DialogHeader>
                        <Button type="submit" onClick={() => window.open(removeApiResponse?.finalUrl)}>Download</Button>
                    </DialogContent>
                </Dialog>
            )
        );
    }

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100 dark:bg-gray-900">
            <div className="w-full max-w-md">
                <Card>
                    <CardHeader>
                        <div className="flex flex-row justify-between">
                        <CardTitle>Video to Article</CardTitle>
                            <Button size="icon" variant="ghost" onClick={handleSettingSubmit}>
                                <SettingsIcon className="h-6 w-6" />
                            </Button>
                        </div>
                        <CardDescription>Upload your video link or file to get the video content</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="video-link">Video Link</Label>
                            <Input id="video-link" value={videoLink} onChange={handleVideoLinkChange}
                                   placeholder="Enter video link"/>
                            <div className="space-x-2 leading-none flex flex-row justify-start">
                                <Checkbox checked={isRemoveWatermark} onClick={handleRemoveWatermarkChange}/>
                                <label
                                    className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                                >
                                    单纯的去除水印
                                </label>
                                <Checkbox checked={isOpenAI} onClick={handleOpenAIChange}/>
                                <label
                                    className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                                >
                                    OPENAI模型
                                </label>
                                <Checkbox checked={isGemini} onClick={handleGminiChange}/>
                                <label
                                    className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                                >
                                    GEMINI模型
                                </label>
                            </div>
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="video-file">Video File</Label>
                            <Input id="video-file" type="file" onChange={handleVideoFileChange}/>
                        </div>
                        <Button className="w-full" type="button" onClick={handleSubmit} disabled={isLoading}>
                            {isLoading ? (
                                <>
                                    <ReloadIcon className="mr-2 h-4 w-4 animate-spin"/>
                                    Please wait
                                </>
                            ) : (
                                '开始转换'
                            )}
                        </Button>
                    </CardContent>
                </Card>
                <Card className="mt-4">
                    <CardHeader>
                        <CardTitle>Result</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="article">Article</Label>
                            <textarea
                                className="w-full h-48 border rounded shadow p-2"
                                id="article"
                                placeholder="Article will be displayed here"
                                readOnly
                                value={videoLinkApiResponse?.content}
                            />
                        </div>
                        <div className="space-y-2">
                            <Button className="w-full" type="button"
                                    onClick={() => window.open(videoLinkApiResponse?.finalUrl)}
                                    disabled={!videoLinkApiResponse?.finalUrl}>
                                Download Video Link
                            </Button>
                        </div>
                    </CardContent>
                </Card>
                {showDialog && <DownloadDialog/>}
                {showSetting&&<DrawerDialogDemo/>}
            </div>
        </div>
    )
}
