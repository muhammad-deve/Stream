import { useState } from "react";
import { useNavigate } from "react-router-dom";
import Header from "@/components/Header";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Heart, DollarSign, ExternalLink } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { useLanguage } from "@/contexts/LanguageContext";

const Donate = () => {
  const navigate = useNavigate();
  const { toast } = useToast();
  const { t } = useLanguage();
  const [amount, setAmount] = useState("");
  const [customAmount, setCustomAmount] = useState("");

  const presetAmounts = ["5", "10", "25", "50", "100"];

  const handleDonate = (method: string) => {
    const donationAmount = amount || customAmount;
    if (!donationAmount || parseFloat(donationAmount) <= 0) {
      toast({
        title: t("donate.invalidAmount"),
        description: t("donate.invalidAmountDesc"),
        variant: "destructive",
      });
      return;
    }

    if (method === "buymeacoffee") {
      window.open("https://buymeacoffee.com/muhammad_deve", "_blank");
    } else if (method === "tirikchilik") {
      toast({
        title: t("donate.thankYouTitle"),
        description: "Tirikchilik.uz integration coming soon!",
      });
    }
  };

  return (
    <div className="min-h-screen bg-background">
      <Header onSearch={(q) => navigate(`/browse?search=${encodeURIComponent(q)}`)} />

      <div className="container mx-auto py-16 px-4">
        <div className="max-w-4xl mx-auto">
          {/* Header */}
          <div className="text-center mb-12">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary/20 mb-6">
              <Heart className="h-8 w-8 text-primary" />
            </div>
            <h1 className="text-4xl font-bold mb-4 text-foreground">{t("donate.title")}</h1>
            <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
              {t("donate.description")}
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            {/* Donation Form */}
            <Card className="p-8 bg-card border-border">
              <h2 className="text-2xl font-bold mb-6 text-foreground">{t("donate.makeTitle")}</h2>

              {/* Preset Amounts */}
              <div className="mb-6">
                <label className="block text-sm font-medium mb-3 text-foreground">
                  {t("donate.selectAmount")}
                </label>
                <div className="grid grid-cols-3 gap-3">
                  {presetAmounts.map((preset) => (
                    <Button
                      key={preset}
                      variant={amount === preset ? "default" : "outline"}
                      onClick={() => {
                        setAmount(preset);
                        setCustomAmount("");
                      }}
                      className={amount === preset ? "bg-primary text-primary-foreground" : ""}
                    >
                      ${preset}
                    </Button>
                  ))}
                </div>
              </div>

              {/* Custom Amount */}
              <div className="mb-6">
                <label className="block text-sm font-medium mb-3 text-foreground">
                  {t("donate.customAmount")}
                </label>
                <div className="relative">
                  <DollarSign className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <Input
                    type="number"
                    placeholder={t("donate.enterAmount")}
                    value={customAmount}
                    onChange={(e) => {
                      setCustomAmount(e.target.value);
                      setAmount("");
                    }}
                    className="pl-9 bg-secondary border-border"
                    min="1"
                    step="1"
                  />
                </div>
              </div>

              {/* Payment Methods */}
              <div className="mb-6">
                <label className="block text-sm font-medium mb-3 text-foreground">
                  {t("donate.paymentMethod")}
                </label>
                <div className="space-y-2">
                  <Button 
                    variant="outline" 
                    className="w-full justify-between" 
                    onClick={() => handleDonate("buymeacoffee")}
                  >
                    <span>Buy Me a Coffee</span>
                    <ExternalLink className="h-4 w-4" />
                  </Button>
                  <Button 
                    variant="outline" 
                    className="w-full justify-between" 
                    onClick={() => handleDonate("tirikchilik")}
                  >
                    <span>Tirikchilik.uz (Test)</span>
                    <ExternalLink className="h-4 w-4" />
                  </Button>
                </div>
              </div>

              <p className="text-xs text-muted-foreground text-center mt-4">
                {t("donate.secure")}
              </p>
            </Card>

            {/* Why Donate */}
            <div className="space-y-6">
              <Card className="p-6 bg-card border-border">
                <h3 className="text-xl font-bold mb-4 text-foreground">{t("donate.whyTitle")}</h3>
                <ul className="space-y-3 text-muted-foreground">
                  <li className="flex items-start gap-3">
                    <Heart className="h-5 w-5 text-primary mt-0.5 flex-shrink-0" />
                    <span>{t("donate.why1")}</span>
                  </li>
                  <li className="flex items-start gap-3">
                    <Heart className="h-5 w-5 text-primary mt-0.5 flex-shrink-0" />
                    <span>{t("donate.why2")}</span>
                  </li>
                  <li className="flex items-start gap-3">
                    <Heart className="h-5 w-5 text-primary mt-0.5 flex-shrink-0" />
                    <span>{t("donate.why3")}</span>
                  </li>
                  <li className="flex items-start gap-3">
                    <Heart className="h-5 w-5 text-primary mt-0.5 flex-shrink-0" />
                    <span>{t("donate.why4")}</span>
                  </li>
                </ul>
              </Card>

              <Card className="p-6 bg-card border-border">
                <h3 className="text-xl font-bold mb-4 text-foreground">{t("donate.promiseTitle")}</h3>
                <p className="text-muted-foreground mb-4">
                  {t("donate.promiseText")}
                </p>
                <p className="text-sm text-muted-foreground">
                  {t("donate.promiseDetail")}
                </p>
              </Card>

              <Card className="p-6 bg-gradient-to-br from-primary/10 to-primary/5 border-primary/20">
                <h3 className="text-xl font-bold mb-2 text-foreground">{t("donate.thankYou")}</h3>
                <p className="text-muted-foreground">
                  {t("donate.thankYouText")}
                </p>
              </Card>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Donate;
