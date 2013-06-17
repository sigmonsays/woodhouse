package main
import (
    "fmt"
    "math/rand"
)
func Greeting(name string) string {
    greetings := []string{
        "Hallo",
        "chíkmàa",
        "tena yistelegn",
        "Selam",
        "Tungjatjeta",
        "Tel nĩdo",
        "صباح الخير",
        "barev",
        "Grüßgott",
        "salam",
        "hello",
        "kaixo",
        "pryvitańnie",
        "namaskar",
        "Wai",
        "koali",
        "Degemer Mad",
        "zdravei",
        "Sua s'dei",
        "hola",
        "sga-noh",
        "hafa adai",
        "moori-bwanj",
        "Shabe Yabebabe Yeshe",
        "你好",
        "mambo",
        "Kia orana",
        "Tansi",
        "dobré ráno",
        "goddag",
        "goedendag",
        "hyvää päivää",
        "salut",
        "dia duit",
        "gamardjoba",
        "Guten Tag",
        "uthegelluthego, h-idiguh-el l-idiguh-o",
        "Γεια σου",
        "Namaste",
        "नमस्ते",
        "góðan dag",
        "السّلام عليكم",
        "ciào",
        "おはよう",
        "Yow Wah gwaan",
    }

    return fmt.Sprintf("%s %s", greetings[rand.Int() % len(greetings)], name)
}
