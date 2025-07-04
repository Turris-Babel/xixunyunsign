//go:build windows
// +build windows
#include <windows.h>
#include <iostream>
#include <string>

// Helper function to convert from the console's ANSI code page to UTF-8
std::string AnsiToUtf8(const char* ansiStr) {
    int wideCharLen = MultiByteToWideChar(CP_ACP, 0, ansiStr, -1, NULL, 0);
    if (wideCharLen == 0) {
        return "";
    }
    wchar_t* wideCharStr = new wchar_t[wideCharLen];
    MultiByteToWideChar(CP_ACP, 0, ansiStr, -1, wideCharStr, wideCharLen);

    int utf8Len = WideCharToMultiByte(CP_UTF8, 0, wideCharStr, -1, NULL, 0, NULL, NULL);
    if (utf8Len == 0) {
        delete[] wideCharStr;
        return "";
    }
    char* utf8Str = new char[utf8Len];
    WideCharToMultiByte(CP_UTF8, 0, wideCharStr, -1, utf8Str, utf8Len, NULL, NULL);

    std::string result(utf8Str);
    delete[] wideCharStr;
    delete[] utf8Str;
    return result;
}

// 定义 DLL 中导出的函数类型
typedef void (*RunCommandFunc)(const char*);

int main(int argc, char* argv[]) {
    // 加载 DLL
    HMODULE hDLL = LoadLibrary(TEXT("xixunyunsign.dll"));
    if (hDLL == NULL) {
        std::cerr << "Error: Could not load xixunyunsign.dll. Error code: " << GetLastError() << std::endl;
        return 1;
    }

    // 获取导出的 RunCommand 函数地址
    RunCommandFunc runCommand = (RunCommandFunc)GetProcAddress(hDLL, "RunCommand");
    if (runCommand == NULL) {
        std::cerr << "Error: Could not find RunCommand function in xixunyunsign.dll. Error code: " << GetLastError() << std::endl;
        FreeLibrary(hDLL);
        return 1;
    }

    // 拼接命令行参数
    std::string args_str;
    for (int i = 1; i < argc; ++i) {
        // Convert argument from ANSI to UTF-8
        std::string arg = AnsiToUtf8(argv[i]);
        // 如果参数包含空格，则用双引号括起来
        if (arg.find(' ') != std::string::npos) {
            args_str += "\"" + arg + "\"";
        } else {
            args_str += arg;
        }
        if (i < argc - 1) {
            args_str += " ";
        }
    }

    std::cout << "Calling RunCommand function from xixunyunsign.dll..." << std::endl;
    // 调用 RunCommand 函数
    runCommand(args_str.c_str());
    std::cout << "RunCommand function call finished." << std::endl;

    // 卸载 DLL
    FreeLibrary(hDLL);
    return 0;
}
