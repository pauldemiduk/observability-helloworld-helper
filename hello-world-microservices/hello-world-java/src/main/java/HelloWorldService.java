import spark.Spark;
import java.util.Optional;

public class HelloWorldService {
    public static void main(String[] args) {
        Spark.port(8080);

        Spark.get("/hello", (req, res) -> {
            String cloudProvider = Optional.ofNullable(System.getenv("CLOUD_PROVIDER")).orElse("unknown");
            String region = Optional.ofNullable(System.getenv("CLOUD_REGION")).orElse("unknown");
            String zone = Optional.ofNullable(System.getenv("CLOUD_ZONE")).orElse("unknown");

            return "Hello, World from Java! Running on " + cloudProvider + " - Region: " + region + ", Zone: " + zone;
        });
    }
}