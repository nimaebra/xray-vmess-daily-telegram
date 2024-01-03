# Xray Vmess Daily Telegram

این روش یک کانفیگ Vmess هست که روزانه تغییر می کند.

</br>

این روش برای تعداد نفرات ۲۰ نفر به پایین مناسب است.
حتما لزومی ندارد که تلگرام داشته باشید و بدون تلگرام هم نصب می شود.

</br>


این روش به یک سرور لینوکس اوبونتو ۲۲ نیاز دارد
همچنین باید ای پی سفید داشته باشید.

</br>


اگر حرفه ای هستید داکیومنت انگیلیسی را مطالعه فرمایید



## نصب

بهتر است پورت ssh رو عوض کنید ولی اگر حوصله ندارید تغییر ندهید

```
echo "Port 9001" >> /etc/ssh/sshd_config
systemctl restart sshd
service ssh restart
```

حالا باید اسکریپت را نصب کنید

```
wget https://raw.githubusercontent.com/majidrezarahnavard/xray-vmess-daily-telegram/main/install.sh
sudo chmod +x /root/xray-configuration/install.sh
bash /root/xray-configuration/install.sh
```


</br>

قسمت بعدی برای تنظیمات هست و می توانید اسم کانفیگ و سایر تنظیمات رو وارد کنید

```
touch /root/xray-configuration/setting.json
echo "{
    \"port\": 443,
    \"bot_token\" : \"\",
    \"chat_id\" : \"\",
    \"dynamic_subscription\" : false,
    \"channel_name\" : \"Sarina_Esmailzadeh\",
    \"send_vnstat\" : false,
    \"aggregate_subscriptions\" : [],
    \"send_configuration\" : \"first\",
    \"send_subscriptions\" : true,
    \"random_header\" : true
}">  /root/xray-configuration/setting.json
```

</br>

حالا باید اسکریپت را اجرا کنید

```
cd /root/xray-configuration
./xray-telegram
```


آدرس زیر را باز کنید

http://your_ip/subscribe.txt

هر شب ساعت ۱۲ شب به وقت تهران کانفیگ ها عوض می شود.

</br>



اگر می خواهید برای کانال تلگرامی بسازید داکیومنت انگیلیسی را بخوانید