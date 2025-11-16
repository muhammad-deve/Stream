import { createContext, useContext, useState, ReactNode } from "react";

export type Language = "en" | "ru" | "uz";

interface LanguageContextType {
  language: Language;
  setLanguage: (lang: Language) => void;
  t: (key: string) => string;
}

const translations: Record<Language, Record<string, string>> = {
  en: {
    // Header
    "header.home": "Home",
    "header.browse": "Browse",
    "header.donate": "Donate",
    "header.search": "Search channels...",
    // Hero
    "hero.badge": "100% Free • No Registration Required",
    "hero.title": "Watch 3500+ TV Channels",
    "hero.titleHighlight": "Absolutely Free",
    "hero.description": "Stream live TV from around the world. News, sports, movies, music, and more. No credit card. No signup.",
    "hero.browseButton": "Browse All Channels",
    "hero.supportButton": "Support Us",
    // Sections
    "section.featured": "Featured Channels",
    "section.browseCategory": "Browse by Category",
    "section.viewAll": "View All Channels",
    // Footer
    "footer.about": "About Streamly",
    "footer.aboutText": "We provide free access to thousands of TV channels from around the world. No registration, no fees - just pure entertainment.",
    "footer.quickLinks": "Quick Links",
    "footer.support": "Support Us",
    "footer.supportText": "If you enjoy our service, please consider supporting us to keep the platform running.",
    "footer.makeDonation": "Make a Donation",
    "footer.copyright": "© 2025 Streamly. All rights reserved. Streaming made simple.",
    "footer.attribution": "All channels are sourced from publicly available streams via",
    // Browse
    "browse.title": "Browse All Channels",
    "browse.searchResults": "Search Results for",
    "browse.showing": "Showing",
    "browse.of": "of",
    "browse.channels": "channels",
    "browse.category": "Category",
    "browse.country": "Country",
    "browse.language": "Language",
    "browse.noChannels": "No channels found",
    "browse.previous": "Previous",
    "browse.next": "Next",
    "browse.page": "Page",
    // Channel Page
    "channel.back": "Back",
    "channel.notFound": "Channel Not Found",
    "channel.goHome": "Go Home",
    "channel.more": "More",
    "channel.channels": "Channels",
    // Donate
    "donate.title": "Support Streamly",
    "donate.description": "Help us keep 3500+ channels free for everyone. Your donation directly supports our infrastructure and allows us to serve millions of users worldwide.",
    "donate.makeTitle": "Make a Donation",
    "donate.selectAmount": "Select Amount",
    "donate.customAmount": "Or Enter Custom Amount",
    "donate.enterAmount": "Enter amount",
    "donate.paymentMethod": "Payment Method",
    "donate.donateNow": "Donate Now",
    "donate.secure": "Your donation is secure and helps keep our service free for everyone.",
    "donate.whyTitle": "Why Your Support Matters",
    "donate.why1": "Keep all 3500+ channels completely free",
    "donate.why2": "Maintain and upgrade our streaming infrastructure",
    "donate.why3": "Add more channels and improve video quality",
    "donate.why4": "Cover bandwidth costs for millions of users",
    "donate.promiseTitle": "Our Promise",
    "donate.promiseText": "Streamly will always remain free and accessible to everyone, regardless of donations. Your contribution helps us serve more users and improve the experience for all.",
    "donate.promiseDetail": "100% of donations go directly to infrastructure costs and platform improvements.",
    "donate.thankYou": "Thank You! ❤️",
    "donate.thankYouText": "Your generosity helps keep free entertainment accessible to millions of people around the world.",
    "donate.invalidAmount": "Invalid Amount",
    "donate.invalidAmountDesc": "Please select or enter a valid donation amount.",
    "donate.thankYouTitle": "Thank You!",
    "donate.thankYouDesc": "Your donation helps keep Streamly running. Redirecting to payment...",
    // Filters
    "filter.all": "All",
  },
  ru: {
    // Header
    "header.home": "Главная",
    "header.browse": "Обзор",
    "header.donate": "Поддержать",
    "header.search": "Поиск каналов...",
    // Hero
    "hero.badge": "100% Бесплатно • Без регистрации • Без рекламы",
    "hero.title": "Смотрите 3500+ ТВ Каналов",
    "hero.titleHighlight": "Совершенно бесплатно",
    "hero.description": "Транслируйте прямой эфир со всего мира. Новости, спорт, фильмы, музыка и многое другое. Без кредитной карты. Без регистрации. Без рекламы.",
    "hero.browseButton": "Просмотреть все каналы",
    "hero.supportButton": "Поддержать нас",
    // Sections
    "section.featured": "Избранные каналы",
    "section.browseCategory": "Обзор по категориям",
    "section.viewAll": "Просмотреть все каналы",
    // Footer
    "footer.about": "О Streamly",
    "footer.aboutText": "Мы предоставляем бесплатный доступ к тысячам телеканалов со всего мира. Без регистрации, без платежей, без рекламы - просто развлечения.",
    "footer.quickLinks": "Быстрые ссылки",
    "footer.support": "Поддержите нас",
    "footer.supportText": "Если вам нравится наш сервис, пожалуйста, поддержите нас, чтобы платформа продолжала работать.",
    "footer.makeDonation": "Сделать пожертвование",
    "footer.copyright": "© 2025 Streamly. Все права защищены. Стриминг стал проще.",
    "footer.attribution": "Все каналы получены из общедоступных потоков через",
    // Browse
    "browse.title": "Просмотреть все каналы",
    "browse.searchResults": "Результаты поиска для",
    "browse.showing": "Показано",
    "browse.of": "из",
    "browse.channels": "каналов",
    "browse.category": "Категория",
    "browse.country": "Страна",
    "browse.language": "Язык",
    "browse.noChannels": "Каналы не найдены",
    "browse.previous": "Назад",
    "browse.next": "Вперед",
    "browse.page": "Страница",
    // Channel Page
    "channel.back": "Назад",
    "channel.notFound": "Канал не найден",
    "channel.goHome": "На главную",
    "channel.more": "Больше",
    "channel.channels": "каналов",
    // Donate
    "donate.title": "Поддержать Streamly",
    "donate.description": "Помогите нам сохранить 3500+ каналов бесплатными для всех. Ваше пожертвование напрямую поддерживает нашу инфраструктуру и позволяет обслуживать миллионы пользователей по всему миру.",
    "donate.makeTitle": "Сделать пожертвование",
    "donate.selectAmount": "Выберите сумму",
    "donate.customAmount": "Или введите свою сумму",
    "donate.enterAmount": "Введите сумму",
    "donate.paymentMethod": "Способ оплаты",
    "donate.donateNow": "Пожертвовать сейчас",
    "donate.secure": "Ваше пожертвование защищено и помогает сохранить наш сервис бесплатным для всех.",
    "donate.whyTitle": "Почему важна ваша поддержка",
    "donate.why1": "Сохранить все 3500+ каналов полностью бесплатными",
    "donate.why2": "Поддерживать и улучшать инфраструктуру стриминга",
    "donate.why3": "Добавлять больше каналов и улучшать качество видео",
    "donate.why4": "Покрывать расходы на пропускную способность для миллионов пользователей",
    "donate.promiseTitle": "Наше обещание",
    "donate.promiseText": "Streamly всегда останется бесплатным и доступным для всех, независимо от пожертвований. Ваш вклад помогает нам обслуживать больше пользователей и улучшать опыт для всех.",
    "donate.promiseDetail": "100% пожертвований идут непосредственно на расходы инфраструктуры и улучшения платформы.",
    "donate.thankYou": "Спасибо! ❤️",
    "donate.thankYouText": "Ваша щедрость помогает сделать бесплатные развлечения доступными для миллионов людей по всему миру.",
    "donate.invalidAmount": "Недействительная сумма",
    "donate.invalidAmountDesc": "Пожалуйста, выберите или введите действительную сумму пожертвования.",
    "donate.thankYouTitle": "Спасибо!",
    "donate.thankYouDesc": "Ваше пожертвование помогает Streamly работать. Перенаправление на оплату...",
    // Filters
    "filter.all": "Все",
  },
  uz: {
    // Header
    "header.home": "Bosh sahifa",
    "header.browse": "Ko'rib chiqish",
    "header.donate": "Qo'llab-quvvatlash",
    "header.search": "Kanallarni qidirish...",
    // Hero
    "hero.badge": "100% Bepul • Ro'yxatdan o'tish talab qilinmaydi • Reklamasiz",
    "hero.title": "3500+ TV Kanallarni Tomosha Qiling",
    "hero.titleHighlight": "Mutlaqo bepul",
    "hero.description": "Butun dunyo bo'ylab jonli TV efirini tomosha qiling. Yangiliklar, sport, filmlar, musiqa va boshqalar. Kredit karta yo'q. Ro'yxatdan o'tish yo'q. Reklama yo'q.",
    "hero.browseButton": "Barcha kanallarni ko'rish",
    "hero.supportButton": "Bizni qo'llab-quvvatlang",
    // Sections
    "section.featured": "Tanlangan kanallar",
    "section.browseCategory": "Kategoriya bo'yicha ko'rish",
    "section.viewAll": "Barcha kanallarni ko'rish",
    // Footer
    "footer.about": "Streamly haqida",
    "footer.aboutText": "Biz butun dunyo bo'ylab minglab TV kanallariga bepul kirishni ta'minlaymiz. Ro'yxatdan o'tish, to'lovlar, reklama yo'q - faqat sof ko'ngilochar.",
    "footer.quickLinks": "Tezkor havolalar",
    "footer.support": "Bizni qo'llab-quvvatlang",
    "footer.supportText": "Agar xizmatimiz sizga yoqsa, platformani ishlab turishga yordam berish uchun bizni qo'llab-quvvatlashni o'ylab ko'ring.",
    "footer.makeDonation": "Xayriya qilish",
    "footer.copyright": "© 2025 Streamly. Barcha huquqlar himoyalangan. Oqimni oddiy qildik.",
    "footer.attribution": "Barcha kanallar ommaviy oqimlardan olindi",
    // Browse
    "browse.title": "Barcha kanallarni ko'rish",
    "browse.searchResults": "Qidiruv natijalari",
    "browse.showing": "Ko'rsatilmoqda",
    "browse.of": "dan",
    "browse.channels": "kanallar",
    "browse.category": "Kategoriya",
    "browse.country": "Mamlakat",
    "browse.language": "Til",
    "browse.noChannels": "Kanallar topilmadi",
    "browse.previous": "Oldingi",
    "browse.next": "Keyingi",
    "browse.page": "Sahifa",
    // Channel Page
    "channel.back": "Orqaga",
    "channel.notFound": "Kanal topilmadi",
    "channel.goHome": "Bosh sahifaga",
    "channel.more": "Ko'proq",
    "channel.channels": "kanallar",
    // Donate
    "donate.title": "Streamly ni qo'llab-quvvatlang",
    "donate.description": "Barcha uchun 3,500+ kanallarni bepul saqlashga yordam bering. Sizning xayriyangiz to'g'ridan-to'g'ri infrastrukturamizni qo'llab-quvvatlaydi va butun dunyo bo'ylab millionlab foydalanuvchilarga xizmat ko'rsatishga imkon beradi.",
    "donate.makeTitle": "Xayriya qilish",
    "donate.selectAmount": "Summani tanlang",
    "donate.customAmount": "Yoki o'z summangizni kiriting",
    "donate.enterAmount": "Summani kiriting",
    "donate.paymentMethod": "To'lov usuli",
    "donate.donateNow": "Hozir xayriya qilish",
    "donate.secure": "Sizning xayriyangiz xavfsiz va xizmatimizni barcha uchun bepul saqlashga yordam beradi.",
    "donate.whyTitle": "Qo'llab-quvvatlashingiz nima uchun muhim",
    "donate.why1": "Barcha 3500+ kanallarni to'liq bepul saqlash",
    "donate.why2": "Oqim infrastrukturasini saqlash va yangilash",
    "donate.why3": "Ko'proq kanallar qo'shish va video sifatini yaxshilash",
    "donate.why4": "Millionlab foydalanuvchilar uchun tarmoqli kengligi xarajatlarini qoplash",
    "donate.promiseTitle": "Bizning va'damiz",
    "donate.promiseText": "Streamly har doim barcha uchun bepul va ochiq bo'lib qoladi, xayriyalardan qat'i nazar. Sizning hissangiz bizga ko'proq foydalanuvchilarga xizmat ko'rsatishga va barcha uchun tajribani yaxshilashga yordam beradi.",
    "donate.promiseDetail": "100% xayriyalar to'g'ridan-to'g'ri infrastruktura xarajatlariga va platformani yaxshilashga ketadi.",
    "donate.thankYou": "Rahmat! ❤️",
    "donate.thankYouText": "Sizning saxiyligingiz butun dunyo bo'ylab millionlab odamlar uchun bepul ko'ngilocharni ochiq qilishga yordam beradi.",
    "donate.invalidAmount": "Noto'g'ri summa",
    "donate.invalidAmountDesc": "Iltimos, to'g'ri xayriya summasini tanlang yoki kiriting.",
    "donate.thankYouTitle": "Rahmat!",
    "donate.thankYouDesc": "Sizning xayriyangiz Streamly ni ishlab turishga yordam beradi. To'lovga yo'naltirish...",
    // Filters
    "filter.all": "Barchasi",
  },
};

const LanguageContext = createContext<LanguageContextType | undefined>(undefined);

export const LanguageProvider = ({ children }: { children: ReactNode }) => {
  const [language, setLanguage] = useState<Language>("en");

  const t = (key: string): string => {
    return translations[language][key] || key;
  };

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  );
};

export const useLanguage = () => {
  const context = useContext(LanguageContext);
  if (context === undefined) {
    throw new Error("useLanguage must be used within a LanguageProvider");
  }
  return context;
};
