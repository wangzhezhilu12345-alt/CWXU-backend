package eva_task

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"edu-evaluation-backed/internal/data/dal"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/signintech/gopdf"
)

// ExportResult 导出结果
type ExportResult struct {
	XlsxPath string   // xlsx 文件路径
	PdfPaths []string // PDF 文件路径列表
	ZipPath  string   // zip 文件路径
}

// generatePDF 生成单个教师的 PDF 报告
func generatePDF(detail dal.TeacherEvaluationDetail, outputDir string) (string, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	fontPath := "./res/NotoSansSC.ttf"
	if err := pdf.AddTTFFont("NotoSansSC", fontPath); err != nil {
		return "", fmt.Errorf("加载字体失败: %w", err)
	}

	pdf.AddPage()

	pageWidth := gopdf.PageSizeA4.W
	pageHeight := gopdf.PageSizeA4.H
	margin := 50.0
	var currentY float64 = 60.0

	getCenterX := func(text string, fontSize int) float64 {
		pdf.SetFontSize(float64(fontSize))
		width, _ := pdf.MeasureTextWidth(text)
		return (pageWidth - width) / 2
	}

	// 辅助函数：绘制带背景色的区块标题
	drawSectionTitle := func(title string, y float64) {
		pdf.SetFillColor(240, 240, 240) // 浅灰色背景
		// 绘制一个覆盖整行的背景矩形
		pdf.RectFromUpperLeftWithStyle(margin, y-2, pageWidth-(margin*2), 20, "F")

		pdf.SetFillColor(0, 0, 0)         // 恢复文字颜色为黑色
		pdf.SetFont("NotoSansSC", "", 13) // 字号稍微加大
		pdf.SetTextColor(40, 40, 40)
		pdf.SetXY(margin+5, y+2) // 文字往右偏移一点，不在背景边缘
		pdf.Cell(nil, title)
	}

	// ========== 1. 机构抬头 ==========
	pdf.SetFont("NotoSansSC", "", 10)
	pdf.SetTextColor(120, 120, 120)
	orgText := "无锡学院 - 新西伯利亚学院 (Wuxi University - Novosibirsk Institute)"
	pdf.SetXY(getCenterX(orgText, 10), currentY)
	pdf.Cell(nil, orgText)

	currentY += 30

	// ========== 2. 主标题 ==========
	pdf.SetFont("NotoSansSC", "", 24)
	pdf.SetTextColor(30, 80, 162)
	title := "教师评教报告"
	pdf.SetXY(getCenterX(title, 24), currentY)
	pdf.Cell(nil, title)

	currentY += 25 // 【修正】加大间距，防止重叠

	pdf.SetFont("NotoSansSC", "", 12)
	subTitle := "Teacher Evaluation Report"
	pdf.SetXY(getCenterX(subTitle, 12), currentY)
	pdf.Cell(nil, subTitle)

	currentY += 25

	// ========== 3. 分隔线 ==========
	pdf.SetLineWidth(1.0)
	pdf.SetStrokeColor(30, 80, 162)
	pdf.Line(margin, currentY, pageWidth-margin, currentY)

	currentY += 40

	// ========== 4. 基本信息 (带背景色副标题) ==========
	drawSectionTitle("基本信息 / Basic Information", currentY)
	currentY += 35 // 留出背景块的高度和间距

	pdf.SetFont("NotoSansSC", "", 11)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(margin, currentY)
	pdf.Cell(nil, fmt.Sprintf("教师姓名 (Teacher): %s    |    工号 (Staff ID): %s", detail.TeacherName, detail.WorkNo))
	currentY += 22
	pdf.SetXY(margin, currentY)
	pdf.Cell(nil, fmt.Sprintf("课程名称 (Course): %s    |    班级名称 (Class): %s", detail.CourseName, detail.ClassName))

	currentY += 45

	// ========== 5. 得分情况 (带背景色副标题) ==========
	drawSectionTitle("得分情况 / Evaluation Scores", currentY)
	currentY += 35

	pdf.SetFont("NotoSansSC", "", 11)
	pdf.SetXY(margin, currentY)
	pdf.Cell(nil, fmt.Sprintf("综合得分 (Global Score): %.2f", detail.AvgScore))
	currentY += 22
	pdf.SetXY(margin, currentY)
	pdf.Cell(nil, fmt.Sprintf("全部教师排名 (Rank): 第 %d 名 / 共 %d 位教师", detail.Rank, detail.TotalTeachers))

	currentY += 45

	// ========== 6. 学生总结 (带背景色副标题) ==========
	if len(detail.Summaries) > 0 {
		drawSectionTitle("学生总结 / Student Comments", currentY)
		currentY += 35

		pdf.SetFont("NotoSansSC", "", 11)
		for i, summary := range detail.Summaries {
			if currentY > pageHeight-100 {
				pdf.AddPage()
				currentY = 60
			}
			pdf.SetXY(margin, currentY)
			text := fmt.Sprintf("%d. %s", i+1, strings.TrimSpace(summary))
			pdf.MultiCell(&gopdf.Rect{W: pageWidth - (margin * 2), H: 18}, text)
			currentY += 25
		}
	}

	// ========== 7. 页脚 ==========
	footerY := pageHeight - 60.0
	pdf.SetLineWidth(0.5)
	pdf.SetStrokeColor(200, 200, 200)
	pdf.Line(margin, footerY, pageWidth-margin, footerY)

	pdf.SetFont("NotoSansSC", "", 9)
	pdf.SetTextColor(150, 150, 150)
	timeStr := fmt.Sprintf("报告生成时间 (Generated at): %s", time.Now().Format("2006-01-02 15:04:05"))
	pdf.SetXY(getCenterX(timeStr, 9), footerY+15)
	pdf.Cell(nil, timeStr)

	// 保存文件逻辑保持不变...
	safeName := func(s string) string {
		return strings.NewReplacer("/", "-", "\\", "-", ":", "-", "*", "-", "?", "-", "\"", "-", "<", "-", ">", "-", "|", "-").Replace(s)
	}
	fileName := fmt.Sprintf("%s-%s-%s.pdf", safeName(detail.TeacherName), safeName(detail.CourseName), safeName(detail.ClassName))
	filePath := filepath.Join(outputDir, fileName)
	if err := pdf.WritePdf(filePath); err != nil {
		return "", err
	}
	return filePath, nil
}

// generateAllPDFs 为所有教师生成 PDF 报告
func generateAllPDFs(details []dal.TeacherEvaluationDetail, outputDir string) []string {
	var pdfPaths []string

	for _, detail := range details {
		path, err := generatePDF(detail, outputDir)
		if err != nil {
			log.Info("生成 PDF 失败: ", err)
			continue
		}
		pdfPaths = append(pdfPaths, path)
	}

	return pdfPaths
}

// zipFiles 将指定文件打包成 zip 文件
// baseDir: 基础目录
// zipPath: 输出的 zip 文件路径
// files: 要打包的文件路径列表（相对于 baseDir）
func zipFiles(baseDir, zipPath string, files ...string) error {
	// 构建 zip 命令，-j 表示只存储文件路径
	args := []string{"-j", zipPath}
	for _, f := range files {
		args = append(args, filepath.Join(baseDir, f))
	}
	cmd := exec.Command("zip", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("zip failed: %v, output: %s", err, string(output))
	}
	return nil
}
