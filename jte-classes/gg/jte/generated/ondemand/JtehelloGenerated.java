package gg.jte.generated.ondemand;
import com.github.jakedegiovanni.comicshelf.HelloController.HelloModel;
public final class JtehelloGenerated {
	public static final String JTE_NAME = "hello.jte";
	public static final int[] JTE_LINE_INFO = {0,0,2,2,2,4,4,4,4,5,5,5,2,2,2,2};
	public static void render(gg.jte.html.HtmlTemplateOutput jteOutput, gg.jte.html.HtmlInterceptor jteHtmlInterceptor, HelloModel model) {
		jteOutput.writeContent("\r\nHello, ");
		jteOutput.setContext("html", null);
		jteOutput.writeUserContent(model.msg());
		jteOutput.writeContent("\r\n");
	}
	public static void renderMap(gg.jte.html.HtmlTemplateOutput jteOutput, gg.jte.html.HtmlInterceptor jteHtmlInterceptor, java.util.Map<String, Object> params) {
		HelloModel model = (HelloModel)params.get("model");
		render(jteOutput, jteHtmlInterceptor, model);
	}
}
