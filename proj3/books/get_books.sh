for ((i=1;i<=1000;i++)); 
do 
   wget https://www.gutenberg.org/files/$i/$i.txt
   wget https://www.gutenberg.org/files/$i/$i-0.txt
done