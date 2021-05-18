import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class loggerTest {
  @Test
  public void main() {
    System.out.println("hello");
    Logger logger = LoggerFactory.getLogger(getClass());
    System.out.println(logger.toString());
    System.out.println(logger.getClass());
    logger.error("world");
  }
}
